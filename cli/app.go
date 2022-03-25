package cli

import (
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

// App is the main entry point for the CLI
type App struct {
	cli.App
	// apiClient
}

// NewCLIApp creates a new instance of CLIApp
func NewCLIApp() App {
	return App{
		cli.App{
			Name:                 "scientia",
			Usage:                "Scientia is a command line interface for the Scientia API",
			EnableBashCompletion: true,
			Commands:             commands,
			ExitErrHandler: func(c *cli.Context, err error) {
				if err != nil {
					log.Error(err)
				}
			},
		},
	}
}
