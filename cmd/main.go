package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"github.com/troykinsella/crash/app"
)

const (
	version = "0.0.1"
)

func defCommands(app *cli.App) {
	app.Commands = []cli.Command{
		*newTestCommand(),
	}
}

func defGlobalFlags(app *cli.App) {

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

func newApp(c *cli.Context) *app.Crash {
	crash, err := app.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return crash
}

func main() {
	app := newCliApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
