package cli

import (
	"fmt"
	"scientia-cli/scientia"
	"syscall"

	cli "github.com/urfave/cli/v2"
	"golang.org/x/term"
)

var commands = []*cli.Command{
	{
		Name:  "login",
		Usage: "login to the scientia API",
		Action: func(c *cli.Context) error {
			username := c.Args().Get(0)
			fmt.Print("Enter password: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return err
			}
			password := string(bytePassword)
			client := scientia.NewAPIClient()
			
			err = client.Login(username, password)
			return err
		},
		// Flags: ,
	},
}
