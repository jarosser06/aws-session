// +build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func promptMFAToken() (string, error) {
	fmt.Printf("MFA Token: ")
	pass, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return string(pass), err
	}

	fmt.Println()
	return string(pass), nil
}

func defaultConfig() string {
	return filepath.Join(
		os.Getenv("HOMEDRIVE"),
		os.Getenv("HOMEPATH"),
		ConfigFilename,
	)
}
