package main

import (
	"os"

	"github.com/goDoCer/scientiaCLI/cli"
)

func main() {
	app := cli.NewCLIApp()
	app.Run(os.Args)
}

// Installing the CLI
// 		1. Linux - update the installer script
// 		2. Windows - create an msi, config.json
// 		3. Mac - installer script?
