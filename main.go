package main

import (
	"os"

	"github.com/goDoCer/scientiaCLI/cli"
)

func main() {
	app := cli.NewCLIApp()
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}
