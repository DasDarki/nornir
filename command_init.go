package main

import (
	"nornir/io"

	"github.com/urfave/cli/v2"
)

func InitCommand() *cli.Command {
	return &cli.Command{
		Name:    "init",
		Usage:   "Initialize a new negroni configuration",
		Aliases: []string{"i"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "restlib",
				Aliases:  []string{"rest", "r"},
				Usage:    "Choose a rest library to use",
				Required: false,
				Value:    "axios",
			},
			&cli.StringFlag{
				Name:     "output-dir",
				Aliases:  []string{"output", "o"},
				Usage:    "Output directory",
				Required: false,
				Value:    "./out",
			},
			&cli.BoolFlag{
				Name:     "overwrite",
				Aliases:  []string{"ow", "force", "f"},
				Usage:    "Overwrite existing configuration",
				Required: false,
				Value:    false,
			},
			&cli.StringFlag{
				Name:     "prefix",
				Aliases:  []string{"p"},
				Usage:    "Prefix for generated files",
				Required: false,
				Value:    "_nornir_",
			},
		},
		Action: onInitCommand,
	}
}

func onInitCommand(ctx *cli.Context) error {
	overwrite := ctx.Bool("overwrite")
	if io.IsConfigExists() && !overwrite {
		return cli.Exit("Negroni configuration already exists. Use --overwrite to overwrite existing configuration", 1)
	}

	restLib := ctx.String("restlib")
	outputDir := ctx.String("output-dir")
	prefix := ctx.String("prefix")

	if restLib != "axios" && restLib != "fetch" {
		return cli.Exit("Invalid rest library (axios, fetch)", 1)
	}

	cfg := &io.Config{
		RestLib:   restLib,
		OutputDir: outputDir,
		Prefix:    prefix,
	}

	cfg.Save()

	return cli.Exit("New negroni configuration created", 0)
}
