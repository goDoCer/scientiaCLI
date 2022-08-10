package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// saveDirCmd represents the saveDir command
var (
	dir        string
	saveDirCmd = &cobra.Command{
		Use:   "save-dir <dir>",
		Short: "set the directory to save downloaded files",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("the save-dir arg is argument is required")
			}
			dir = args[0]
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(dir)
			absPath, err := filepath.Abs(dir)
			if err != nil {
				return err
			}

			if _, err := os.Stat(absPath); os.IsNotExist(err) {
				log.Warnf("Directory %s does not exist", absPath)
				// TODO: make creating interactive (prompt and read input)
				log.Println("Trying to create directory", absPath)
				log.Println("Should I create the directory?", absPath)
				err = os.Mkdir(absPath, 0777)
				if err != nil {
					log.Fatal("Could not create directory", absPath, " because of ", err)
					return err
				}
			}

			cfg.SaveDir = absPath
			log.Info("Save directory set to ", cfg.SaveDir)

			return cfg.save(configPath)
		},
	}
)

func init() {
	rootCmd.AddCommand(saveDirCmd)
}
