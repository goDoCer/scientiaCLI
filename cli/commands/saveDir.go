package commands

import (
	"github.com/goDoCer/scientiaCLI/cli/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var SaveDir = &cli.Command{
	Name:  "save-dir",
	Usage: "set the directory to save downloaded files",
	Action: func(c *cli.Context) error {
		dir := c.Args().First()
		if dir == "" {
			return errors.New("Please enter a directory")
		}
		config.Cfg.SaveDir = dir
		log.Info("Save directory set to ", config.Cfg.SaveDir)
		return config.Cfg.Save(config.Path)
	},
}
