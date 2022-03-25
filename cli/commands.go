package cli

import (
	"fmt"
	"os"
	"path"
	"scientia-cli/scientia"
	"syscall"

	cli "github.com/urfave/cli/v2"
	"golang.org/x/term"
)

var client scientia.APIClient
var tokenPath string

func init() {
	client = scientia.NewAPIClient()
	filepath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	tokenPath = path.Dir(filepath) + "/token.txt"
	loginDetails, _ := loadDetails()
	client.AddTokens(*loginDetails)
}

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
			err = client.Login(username, password)
			if err != nil {
				return err
			}
			return saveDetails(client.GetTokens())
		},
		// Flags: ,
	},
	{
		Name:  "download",
		Usage: "download a file from scientia",
		Action: func(c *cli.Context) error {
			client.GetCourses()
			return nil
		},
	},
}
