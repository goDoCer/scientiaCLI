package commands

import (
	"fmt"
	"github.com/goDoCer/scientiaCLI/cli/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
	"syscall"
)

var Login = &cli.Command{
	Name:      "login",
	Usage:     "login to the scientia API",
	ArgsUsage: "shortcode",
	Action: func(c *cli.Context) error {
		var shortcode string

		if c.NArg() < 1 {
			log.Warn("No shortcode provided")
			fmt.Print("Please enter your shortcode: ")
			_, err := fmt.Scanln(&shortcode)
			if err != nil {
				return err
			}
		} else {
			shortcode = c.Args().First()
			log.Infof("shortcode: %s", shortcode)
		}

		if shortcode == "" {
			return errors.New("Please enter your shortcode")
		}

		fmt.Print("Enter password (It will be hidden): ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Println()

		password := string(bytePassword)
		err = config.Client.Login(shortcode, password)
		if err != nil {
			return err
		}

		config.Cfg.UpdateTokens(config.Client.GetTokens())
		err = config.Cfg.Save(config.Path)
		if err == nil {
			log.Info("Login successful")
		}
		return err
	},
}
