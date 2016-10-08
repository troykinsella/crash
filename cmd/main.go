package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

const (
	version = "0.0.3"

	nocolor = "nc"
)

func defCommands(app *cli.App) {
	app.Commands = []cli.Command{
		*newTestCommand(),
	}
}

func defGlobalFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: nocolor,
			Usage: "No color. Disable output colorization.",
		},
	}
}

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "crash"
	app.Version = version
	app.Usage = "test utility"

	defCommands(app)
	defGlobalFlags(app)

	return app
}

func main() {
	app := newCliApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
