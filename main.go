package main

import (
	"os"
	"scientia-cli/cli"
)

func main() {

	app := cli.NewCLIApp()
	app.Run(os.Args)
}
