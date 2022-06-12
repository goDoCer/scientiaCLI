package cli

import (
	"os"
	"path"

	"github.com/goDoCer/scientiaCLI/logging"
	"github.com/goDoCer/scientiaCLI/scientia"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

var (
	client     scientia.APIClient
	cfg        config
	configPath string
	verbose    bool

	errNotLoggedIn = errors.New("not logged in, login using the login command")
)

// App is the main entry point for the CLI
type App struct {
	cli.App
}

// NewCLIApp creates a new instance of CLIApp
func NewCLIApp() App {
	return App{
		App: cli.App{
			Name:                 "scientia",
			Usage:                "Scientia is a command line interface for the Scientia API",
			EnableBashCompletion: true,
			Commands:             commands,
			ExitErrHandler: func(c *cli.Context, err error) {
				if err != nil {
					log.Error(err)
				}
			},

			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "verbose",
					Aliases:     []string{"v"},
					Usage:       "Enable logging",
					Value:       false,
					Destination: &verbose,
				},
			},

			Before: func(ctx *cli.Context) error {
				client = scientia.NewAPIClient()

				filepath, err := os.Executable()
				if err != nil {
					return err
				}

				// reading the config from the executable's directory
				configPath = path.Join(path.Dir(filepath), "config.json")
				cfg, err = loadConfig(configPath)
				if err != nil {
					return err
				}
				tokens, found := cfg.tokens()
				if !found {
					return errNotLoggedIn
				}

				client.AddTokens(tokens)

				log.SetOutput(logging.L)

				return nil
			},
			After: func(ctx *cli.Context) error {
				tokens := client.GetTokens()
				cfg.updateTokens(tokens)

				if verbose {
					logging.L.Print()
				}
				return nil
			},
		},
	}
}
