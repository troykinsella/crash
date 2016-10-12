package main

import (
	"github.com/urfave/cli"
	"github.com/troykinsella/crash/app"
)

func newValidateCommand() *cli.Command {
	return &cli.Command{
		Name:        "validate",
		Aliases:     []string{"v"},
		Usage:       "Validate a Crashfile test plan without running it",
		Action: func(c *cli.Context) error {
			a := app.New()
			err := a.Validate(c.String(testFile))
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			return nil
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  testFile,
				Usage: "Crashfile test plan yaml `FILE`",
			},
		},
	}
}
