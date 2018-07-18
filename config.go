package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

type Alias struct {
	AccountNumber int    `yaml:"account_number" required:"true"`
	DefaultRegion string `yaml:"default_region"`
	Name          string `yaml:"name" required:"true"`
	Role          string `yaml:"role" required:"true"`
}

type Config struct {
	Aliases            []Alias `yaml:"aliases"`
	AWSAccessKeyId     string  `yaml:"aws_access_key_id" required:"true"`
	AWSSecretAccessKey string  `yaml:"aws_secret_access_key" required:"true"`
	MFARole            string  `yaml:"mfa_role" required:"true"`
	aliasMap           map[string]int
}

// Return an Alias based on the name
func (c *Config) GetAlias(name string) (*Alias, error) {
	if alias, ok := c.aliasMap[name]; ok {
		return &c.Aliases[alias], nil
	}

	return nil, fmt.Errorf("alias %s does not exist", name)
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
	config.aliasMap = make(map[string]int)
	for index, alias := range config.Aliases {
		config.aliasMap[alias.Name] = index
	}

	return &config, nil
}
