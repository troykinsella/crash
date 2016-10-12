package main

import (
	"github.com/troykinsella/crash"
	"errors"
	"github.com/urfave/cli"
	"strings"
	"github.com/troykinsella/crash/app"
	"github.com/troykinsella/crash/logging"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

const (
	json        = "j"
	quiet       = "q"
	verbose1    = "v"
	verbose2    = "vv"
	verbose3    = "vvv"
	setVariable = "s"
	testFile    = "f"
	nocolor     = "nc"
)

func loadVariablesFile(file string, vars map[string]string) error {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(fileBytes, vars)
	if err != nil {
		return err
	}
	return nil
}

func loadVariables(vals []string, vars map[string]string) error {
	for _, val := range vals {
		parts := strings.Split(val, "=")
		switch len(parts) {
		case 1:
			if err := loadVariablesFile(val, vars); err != nil {
				return err
			}
		case 2:
			vars[parts[0]] = parts[1]
		default:
			return errors.New("Malformed varible: " + val)
		}
	}
	return nil
}

func logLevel(c *cli.Context) (logging.Level) {
	if c.Bool(quiet) {
		return logging.L_OFF
	}
	if c.Bool(verbose1) {
		return logging.L_INFO
	}
	if c.Bool(verbose2) {
		return logging.L_DEBUG
	}
	if c.BoolT(verbose3) {
		return logging.L_TRACE
	}
	return logging.L_DEFAULT
}

func newTestOptions(c *cli.Context) (*crash.TestOptions, error) {
	vars := make(map[string]string)

	err := loadVariables(c.StringSlice(setVariable), vars)
	if err != nil {
		return nil, err
	}

	options := &crash.TestOptions{
		Crashfile: c.String(testFile),
		LogLevel: logLevel(c),
		Colorize: !c.IsSet(nocolor),
		LogJson: c.Bool(json),
		Variables: vars,
	}
	return options, nil
}

func newTestCommand() *cli.Command {
	return &cli.Command{
		Name:        "test",
		Aliases:     []string{"t"},
		Usage:       "Execute a Crashfile test plan",
		Description: "A Crashfile to execute can be specified with the -f option,\n" +
		"   or if omitted, crash searches the current directory for the first match, in order:\n" +
		"     * Crashfile\n" +
		"     * Crashfile.yml\n" +
		"     * Crashfile.yaml\n",
		Action: func(c *cli.Context) error {
			a := app.New()

			opts, err := newTestOptions(c)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			ok, err := a.Test(opts)
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
				Name: json,
				Usage: "Format logging output as JSON",
			},
			cli.BoolFlag{
				Name: nocolor,
				Usage: "No color. Disable output colorization.",
			},
			cli.BoolFlag{
				Name: quiet,
				Usage: "Quiet mode; suppress logging",
			},
			cli.StringSliceFlag{
				Name:  setVariable,
				Usage: "Set variable(s) `FILE|KEY=VALUE`",
			},
			cli.BoolFlag{
				Name: verbose1,
				Usage: "Verbose logging; Use -vv or -vvv to increase verbosity",
			},
			cli.BoolFlag{
				Name: verbose2,
				Hidden: true,
			},
			cli.BoolFlag{
				Name: verbose3,
				Hidden: true,
			},
			cli.StringFlag{
				Name:  testFile,
				Usage: "Crashfile test plan yaml `FILE`; Defaults to searching for Crashfile.y[a]ml in the current directory",
			},
		},
	}
}
