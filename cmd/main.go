package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

const (
	version = "0.0.5"
)

func defCommands(app *cli.App) {
	app.Commands = []cli.Command{
		*newTestCommand(),
		*newValidateCommand(),
	}
}

func defGlobalFlags(app *cli.App) {}

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "crash"
	app.Version = version
	app.Usage = "Run functional and performance tests, asserting and reporting on results.\n" +
	"   For more detailed help, run: " + app.Name + " help <command>"
	app.Author = "Troy Kinsella"

	defCommands(app)
	defGlobalFlags(app)

	return app
}

func main() {
	cli.VersionFlag = cli.BoolFlag{
		Name: "V, version",
		Usage: "print the version",
	}

	app := newCliApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
