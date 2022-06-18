package cli

import (
	cmds "github.com/goDoCer/scientiaCLI/cli/commands"
	"github.com/goDoCer/scientiaCLI/cli/config"
	"os"
	"path"

	"github.com/goDoCer/scientiaCLI/scientia"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

var (
	client     scientia.APIClient
	cfg        config.Config
	configPath string
)

// App is the main entry point for the CLI
type App struct {
	cli.App
}

// NewCLIApp creates a new instance of CLIApp
func NewCLIApp() App {

	// List of all commands available in the CLI
	var allCommands = []*cli.Command{
		cmds.Login,
		cmds.Download,
		cmds.SaveDir,
		cmds.Config,
	}

	return App{
		App: cli.App{
			Name:                 "scientia",
			Usage:                "Scientia is a command line interface for the Scientia API",
			EnableBashCompletion: true,
			Commands:             allCommands,
			ExitErrHandler: func(c *cli.Context, err error) {
				if err != nil {
					log.Error(err)
				}
			},

			Before: func(ctx *cli.Context) error {
				client = scientia.NewAPIClient()

				filepath, err := os.Executable()
				if err != nil {
					return err
				}

				// reading the Config from the executable's directory
				configPath = path.Join(path.Dir(filepath), "config.json")
				cfg, err = config.LoadConfig(configPath)
				if err != nil {
					return err
				}

				tokens := cfg.Tokens()

				client.AddTokens(tokens)

				return nil
			},
			After: func(ctx *cli.Context) error {
				tokens := client.GetTokens()
				cfg.UpdateTokens(tokens)
				return nil
			},
		},
	}
}
