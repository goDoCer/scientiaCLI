package main

import (
	"os"
	"scientia-cli/cli"

	"github.com/sirupsen/logrus"
)

func main() {

	config, err := cli.ReadConfig("./config.json")
	if err != nil {
		logrus.Fatal(err)
	}
	app := cli.NewCLIApp(config)
	app.Run(os.Args)
}
