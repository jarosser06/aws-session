package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"syscall"

	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

func promptMFAToken() (string, error) {
	// Using dev tty
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return "", err
	}

	fmt.Fprintf(tty, "MFA Token: ")
	pass, err := terminal.ReadPassword(int(tty.Fd()))
	if err != nil {
		return string(pass), err
	}

	fmt.Fprintln(tty)
	return string(pass), nil
}

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
	if mfaTok == "" {
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
	if mfaTok == "" {
		mfaTok, err = promptMFAToken()
		if err != nil {
			return err
		}
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
		Region:             c.String("region"),
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
	defaultConfigPath := path.Join(os.Getenv("HOME"), ".tok/config")

	app := cli.NewApp()
	app.Name = "aws-session"
	app.Version = "0.1.0"
	app.Usage = "Provides an easy way to assume roles"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  defaultConfigPath,
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
