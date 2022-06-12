package main

import (
	"os"

	"github.com/goDoCer/scientiaCLI/cli"
)

func main() {
	app := cli.NewCLIApp()
	app.Run(os.Args)
}
