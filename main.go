package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

var Version = ""

func webCommand(c *cli.Context) error {
	config, err := LoadConfig(c.GlobalString("config"))
	if err != nil {
		return err
	}

	aliasName := c.String("alias")
	if aliasName == "" {
		return fmt.Errorf("alias flag can not be empty")
	}

	alias, credentials, err := config.GetAlias(aliasName)
	if err != nil {
		return err
	}

	mfaTok := c.String("token-code")
	if mfaTok == "" && credentials.MFARole != "" {
		mfaTok, err = promptMFAToken()
		if err != nil {
			return err
		}
	}

	input := webOutInput{
		AWSAccessKeyID:     credentials.AWSAccessKeyId,
		AWSSecretAccessKey: credentials.AWSSecretAccessKey,
		AccountName:        alias.Name,
		AWSAccountNumber:   strconv.Itoa(alias.AccountNumber),
		RoleName:           alias.Role,
		MFADeviceID:        credentials.MFARole,
		TokenCode:          mfaTok,
		SessionName:        c.String("session-name"),
		Duration:           c.Int("duration"),
	}

	out, err := webOut(input)
	if err != nil {
		return err
	}

	fmt.Println(out)

	return nil
}

func authCommand(c *cli.Context) error {
	config, err := LoadConfig(c.GlobalString("config"))
	if err != nil {
		return err
	}

	aliasName := c.String("alias")
	if aliasName == "" {
		return fmt.Errorf("alias flag can not be empty")
	}

	alias, credentials, err := config.GetAlias(aliasName)
	if err != nil {
		return err
	}

	mfaTok := c.String("token-code")
	if mfaTok == "" && credentials.MFARole != "" {
		mfaTok, err = promptMFAToken()
		if err != nil {
			return err
		}
	}

	region := c.String("region")
	if alias.DefaultRegion != "" {
		region = alias.DefaultRegion
	}

	input := credentialsOutInput{
		AWSAccessKeyID:     credentials.AWSAccessKeyId,
		AWSSecretAccessKey: credentials.AWSSecretAccessKey,
		AccountName:        alias.Name,
		AWSAccountNumber:   strconv.Itoa(alias.AccountNumber),
		RoleName:           alias.Role,
		MFADeviceID:        credentials.MFARole,
		TokenCode:          mfaTok,
		SessionName:        c.String("session-name"),
		Duration:           c.Int("duration"),
		Region:             region,
		UserShell:          c.String("format"),
	}

	out, err := credentialOut(input)
	if err != nil {
		return err
	}

	fmt.Println(out)

	return nil
}

func listCommand(c *cli.Context) error {
	config, err := LoadConfig(c.GlobalString("config"))
	if err != nil {
		return err
	}

	for _, name := range config.AliasNames() {
		fmt.Println(name)
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "aws-session"
	app.Version = Version
	app.Usage = "Provides an easy way to assume roles"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  defaultConfig(),
			Usage:  "Tok config",
			EnvVar: "TOK_CONFIG",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "list",
			Usage:  "List available aliases",
			Action: listCommand,
		},
		{
			Name:  "auth",
			Usage: "Get credentials",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "alias, A",
					Value:  "",
					Usage:  "Account Alias to fetch credentials for",
					EnvVar: "TOK_ALIAS",
				},
				cli.StringFlag{
					Name:  "format, F",
					Value: "",
					Usage: "Format to expose secrets in Must be one of powershell, cmd, docker, bash.",
				},
				cli.StringFlag{
					Name:   "token-code, T",
					Value:  "",
					Usage:  "MFA Token",
					EnvVar: "TOK_TOKEN",
				},
				cli.StringFlag{
					Name:   "region, r",
					Value:  "",
					Usage:  "AWS Region to include in ENV Variables.",
					EnvVar: "TOK_REGION",
				},
				cli.StringFlag{
					Name:  "session-name n",
					Value: "",
					Usage: "Optional session name, will be generated if not set",
				},
				cli.IntFlag{
					Name:  "duration, d",
					Value: 3600,
					Usage: "Credential duration, default 3600",
				},
			},
			Action: authCommand,
		},
		{
			Name:  "web",
			Usage: "Generate Console Signin URL",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "alias, A",
					Value:  "",
					Usage:  "Account Alias to fetch credentials for",
					EnvVar: "TOK_ALIAS",
				},
				cli.StringFlag{
					Name:   "token-code, T",
					Value:  "",
					Usage:  "MFA Token",
					EnvVar: "TOK_TOKEN",
				},
				cli.StringFlag{
					Name:  "session-name n",
					Value: "",
					Usage: "Optional session name, will be generated if not set",
				},
				cli.IntFlag{
					Name:  "duration, d",
					Value: 3600,
					Usage: "Credential duration, default 3600",
				},
			},
			Action: webCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
