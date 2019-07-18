// +build linux darwin

package main

import (
	"fmt"
	"os"
	"path"
	"syscall"

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

func defaultConfig() string {
	return path.Join(os.Getenv("HOME"), ConfigFilename)
}
