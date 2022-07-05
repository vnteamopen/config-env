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

type FlagName string

const (
	FlagOverwrite      FlagName = "overwrite"
	FlagOutputToScreen FlagName = "out-screen"
	FlagCustomPattern  FlagName = "custom"
)

var Flags = []cli.Flag{
	&cli.BoolFlag{
		Name:     string(FlagOverwrite),
		Usage:    "-w",
		Required: false,
		Value:    false,
		Aliases:  []string{"w"},
	},
	&cli.BoolFlag{
		Name:     string(FlagOutputToScreen),
		Usage:    "-out-screen",
		Required: false,
		Value:    false,
	},
	&cli.StringSliceFlag{
		Name:       string(FlagCustomPattern),
		Usage:      "-c",
		Required:   false,
		HasBeenSet: true,
		Value:      cli.NewStringSlice("{{", "}}"),
		Aliases:    []string{"c"},
	},
}

func main() {
	app := &cli.App{
		Name:     name,
		Version:  version,
		Compiled: time.Now(),
		Authors:  []*cli.Author{&cli.Author{Name: "https://vnteamopen.com"}},
		HelpName: "config-env",
		Usage:    "A tool for embedding file's contents into the input file. Embedded pattern is {{env \"variable_environment_name\"}}",
		UsageText: `config-env /path/to/input/file /path/to/output/file
config-env -c begin-pattern,end-pattern /path/to/input/file /path/to/output/file
config-env -w /path/to/input/file
config-env help`,
		EnableBashCompletion: true,
		Flags:                Flags,
		Action:               Action,
	}
	app.Run(os.Args)
}

func Action(c *cli.Context) error {
	c.App.Setup()

	isOverwrite := c.Bool(string(FlagOverwrite))
	isOutputToScreen := c.Bool(string(FlagOutputToScreen))
	if valid := validArgs(c.NArg(), isOverwrite, isOutputToScreen); !valid {
		cli.ShowAppHelp(c)
		return cli.Exit("", 1)
	}
	templatePath, outputPaths := getPaths(c.Args(), isOverwrite)

	pattern := c.StringSlice(string(FlagCustomPattern))
	if valid := validPattern(pattern); !valid {
		cli.ShowAppHelp(c)
		return cli.Exit("Custom pattern must contain begin and end part", 1)
	}

	if err := actions.CharByCharParse(actions.ParseRequest{
		InputPath:        templatePath,
		ListOutputPath:   outputPaths,
		IsOutputToScreen: isOutputToScreen,
		Pattern:          pattern,
	}); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	if isOverwrite {
		if err := actions.OverwriteInput(templatePath); err != nil {
			return cli.Exit(err.Error(), 1)
		}
	}

	return cli.Exit("", 0)
}

func validArgs(totalArgs int, isOverwrite, isOutputToScreen bool) bool {
	requiredArgs := 2
	if isOverwrite || isOutputToScreen {
		requiredArgs = 1
	}

	return totalArgs >= requiredArgs
}

func getPaths(args cli.Args, isOverwrite bool) (templatePath string, listOutputPath []string) {
	templatePath = args.Get(0)
	noOutputs := args.Len() - 1

	listOutputPath = make([]string, 0, args.Len())
	if isOverwrite {
		listOutputPath = append(listOutputPath, actions.CreateTmpFile(templatePath))
	}
	for i := 0; i < noOutputs; i++ {
		listOutputPath = append(listOutputPath, args.Get(i+1))
	}
	return templatePath, listOutputPath
}

func validPattern(pattern []string) bool {
	if len(pattern) != 2 {
		return false
	}
	return true
}
