package commands

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

var Ls = &cli.Command{
	Name:  "ls",
	Usage: "list the files in the current directory",
	Action: func(c *cli.Context) error {
		fmt.Println("ls command")
		return nil
	},
}
