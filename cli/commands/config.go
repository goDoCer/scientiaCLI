package commands

import (
	"fmt"
	"github.com/goDoCer/scientiaCLI/cli/config"
	"github.com/urfave/cli/v2"
)

var Config = &cli.Command{
	Name:  "config",
	Usage: "view the CLI configuration",
	Action: func(c *cli.Context) error {
		fmt.Println("Save directory: ", config.Cfg.SaveDir)
		fmt.Println("Access token: ", config.Cfg.AccessToken)
		fmt.Println("Refresh token: ", config.Cfg.RefreshToken)
		return nil
	},
}
