package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/goDoCer/scientiaCLI/scientia"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	configPath string
	showConfig bool

	cfg    config
	client = scientia.NewAPIClient()

	rootCmd = &cobra.Command{
		Use:   "scientia-cli",
		Short: "scientia-cli is a command line interface for Scientia",
		// Long: ``, TODO: add long description
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			client = scientia.NewAPIClient()

			filepath, err := os.Executable()
			if err != nil {
				return err
			}

			if configPath == "" {
				// reading the config from the executable's directory
				configPath = path.Join(path.Dir(filepath), "config.json")
			}

			cfg, err = loadConfig(configPath)
			if err != nil {
				return err
			}

			tokens := cfg.tokens()

			client.AddTokens(tokens)

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if showConfig {
				fmt.Println("Save directory: ", cfg.SaveDir)
				fmt.Println("Access token: ", cfg.AccessToken)
				fmt.Println("Refresh token: ", cfg.RefreshToken)
			} else {
				cmd.Help()
			}
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			tokens := client.GetTokens()
			cfg.updateTokens(tokens)
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "C", "", "path to config file")
	rootCmd.PersistentFlags().BoolVarP(&showConfig, "show-config", "s", false, "show-config")

	rootCmd.PersistentFlags().MarkHidden("config")
	rootCmd.PersistentFlags().MarkHidden("show-config")
}
