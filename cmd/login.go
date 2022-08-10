/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to scientia",
	Long:  `TODO login desc`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var shortcode string
		fmt.Print("Please enter your shortcode: ")
		fmt.Scanln(&shortcode)

		fmt.Print("Enter password (It will be hidden): ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Println()

		password := string(bytePassword)
		err = client.Login(shortcode, password)
		if err != nil {
			return err
		}

		cfg.updateTokens(client.GetTokens())
		err = cfg.save(configPath)
		if err == nil {
			log.Info("Login successful")
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
