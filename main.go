package main

import (
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/vnteamopen/config-env/actions"
)

const (
	//  Help template: cli.AppHelpTemplate
	name    = "Config environments"
	version = "1.0.0"
)

func main() {
	app := &cli.App{
		Name:     name,
		Version:  version,
		Compiled: time.Now(),
		Authors:  []*cli.Author{&cli.Author{Name: "https://vnteamopen.com"}},
		HelpName: "config-env",
		Usage:    "A tool for embedding file's contents into the input file. Embedded pattern is {{env \"variable_environment_name\"}}",
		UsageText: `config-env /path/to/input/file /path/to/output/file
config-env help`,
		EnableBashCompletion: true,
		Action:               Action,
	}
	app.Run(os.Args)
}

func Action(c *cli.Context) error {
	c.App.Setup()
	if c.NArg() <= 1 {
		cli.ShowAppHelp(c)
		return cli.Exit("", 0)
	}
	inputPath := c.Args().Get(0)
	outputPath := c.Args().Get(1)
	if err := actions.Parse(inputPath, outputPath); err != nil {
		return cli.Exit(err.Error(), 1)
	}
	return cli.Exit("", 0)
}
