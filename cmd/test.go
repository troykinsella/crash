package main

import (
	"github.com/troykinsella/crash"
	"errors"
	"github.com/urfave/cli"
	"strings"
)

const (
	debug       = "d"
	json        = "j"
	quiet       = "q"
	setVariable = "s"
	testFile    = "f"
	varsFile    = "v"
)

func parseVariables(vals []string, vars *map[string]string) error {
	for _, val := range vals {
		parts := strings.Split(val, "=")
		if len(parts) != 2 {
			return errors.New("Malformed varible: " + val)
		}
		(*vars)[parts[0]] = parts[1]
	}
	return nil
}

func newTestOptions(c *cli.Context) (*crash.TestOptions, error) {
	vars := make(map[string]string)

	err := parseVariables(c.StringSlice(setVariable), &vars)
	if err != nil {
		return nil, err
	}

	options := &crash.TestOptions{
		InputFile: c.String(testFile),
		Debug: c.Bool(debug),
		Quiet: c.Bool(quiet),
		LogJson: c.Bool(json),
		Variables: vars,
		VariablesFiles: c.StringSlice(varsFile),
	}
	return options, nil
}

func newTestCommand() *cli.Command {
	return &cli.Command{
		Name:    "test",
		Aliases: []string{"t"},
		Usage:   "Execute a test plan",
		Action: func(c *cli.Context) error {
			app := newApp(c)

			opts, err := newTestOptions(c)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			ok, err := app.Test(opts)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			if !ok {
				return cli.NewExitError("", 2)
			}

			return nil
		},
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: debug,
				Usage: "Debug mode; increase logging verbosity",
			},
			cli.BoolFlag{
				Name: json,
				Usage: "Format logging output as JSON",
			},
			cli.BoolFlag{
				Name: quiet,
				Usage: "Quiet mode; suppress logging",
			},
			cli.StringSliceFlag{
				Name:  setVariable,
				Usage: "Variable key pair `key=value`",
			},
			cli.StringFlag{
				Name:  testFile,
				Usage: "Input test yaml `FILE`",
			},
			cli.StringSliceFlag{
				Name: varsFile,
				Usage: "Variables yaml `FILE`",
			},
		},
	}
}
