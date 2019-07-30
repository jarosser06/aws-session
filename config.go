package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

const ConfigFilename = ".aws-session/config.yaml"

type aliasLocation struct {
	accountIndex int
	aliasIndex   int
}

type Alias struct {
	AccountNumber int    `yaml:"account_number" required:"true"`
	DefaultRegion string `yaml:"default_region"`
	Name          string `yaml:"name" required:"true"`
	Role          string `yaml:"role" required:"true"`
}

type Account struct {
	Aliases            []Alias `yaml:"aliases"`
	AWSAccessKeyId     string  `yaml:"aws_access_key_id" required:"true"`
	AWSSecretAccessKey string  `yaml:"aws_secret_access_key" required:"true"`
	MFARole            string  `yaml:"mfa_role" required:"true"`
}

type Config struct {
	Accounts []Account `yaml:"accounts"`
	aliasMap map[string]aliasLocation
}

type SecurityCredentials struct {
	AWSAccessKeyId     string
	AWSSecretAccessKey string
	MFARole            string
}

// Return an Alias based on the name
func (c *Config) GetAlias(name string) (*Alias, *SecurityCredentials, error) {
	if location, ok := c.aliasMap[name]; ok {
		account := c.Accounts[location.accountIndex]

		credentials := &SecurityCredentials{
			AWSAccessKeyId:     account.AWSAccessKeyId,
			AWSSecretAccessKey: account.AWSSecretAccessKey,
			MFARole:            account.MFARole,
		}

		alias := &account.Aliases[location.aliasIndex]

		return alias, credentials, nil
	}

	return nil, nil, fmt.Errorf("alias %s does not exist", name)
}

// Return a slice of Alias names
func (c *Config) AliasNames() []string {
	aliases := make([]string, len(c.aliasMap))

	i := 0
	for name := range c.aliasMap {
		aliases[i] = name
		i++
	}

	return aliases
}

// Load Configuration File
func LoadConfig(filePath string) (*Config, error) {
	var config Config

	openFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer openFile.Close()

	b, err := ioutil.ReadAll(openFile)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}

	// Validate Struct
	if err := Validate(config); err != nil {
		return nil, err
	}

	// Populate aliasMap
	config.aliasMap = make(map[string]aliasLocation)
	for accountIndex, account := range config.Accounts {
		for aliasIndex, alias := range account.Aliases {
			config.aliasMap[alias.Name] = aliasLocation{
				accountIndex: accountIndex,
				aliasIndex:   aliasIndex,
			}
		}
	}

	return &config, nil
}
