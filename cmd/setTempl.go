/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"os"
	"text/template"

	"github.com/goDoCer/scientiaCLI/scientia"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// setTemplCmd represents the setTempl command
var setTemplCmd = &cobra.Command{
	Use:   "set-templ [template]",
	Short: "This command sets the Go template for the course directories.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		templStr := args[0]

		templ, err := template.New("template").Parse(templStr)
		if err != nil {
			return errors.New("invalid template string passed in. You can read more about the template syntax here: https://golang.org/pkg/text/template/")
		}

		log.Info("This is what the template will look like for a test course:")
		err = templ.Execute(os.Stdout, scientia.Course{
			Title:        "TestTitle",
			Code:         "TestCode",
			CanManage:    false,
			HasMaterials: false,
		})

		if err != nil {
			return err
		}

		cfg.Template = args[0]
		cfg.save(configPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setTemplCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setTemplCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setTemplCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
