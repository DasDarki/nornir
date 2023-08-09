package main

import (
	"nornir/analyzer"
	"nornir/generator"
	"nornir/io"
	"nornir/log"

	"github.com/urfave/cli/v2"
)

func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:    "generate",
		Usage:   "Generates the project",
		Aliases: []string{"g"},
		Action:  onGenerateCommand,
	}
}

func onGenerateCommand(ctx *cli.Context) error {
	cfg := io.LoadConfig()

	log.Debug("Generating project...")

	rootDir := io.FindFiles()
	modName := io.ReadModuleName()
	log.Debugf("Module name: %s", modName)

	analyzer := analyzer.NewAnalyzer(modName, rootDir, cfg.OutputDir, cfg.Prefix, cfg.RestLib)

	err := analyzer.AnalyzeStructs()
	if err != nil {
		log.Debug("Failed to find structs!")
		return err
	}

	err = analyzer.AnalyzeAnnotations()
	if err != nil {
		log.Debug("Failed to find annotations!")
		return err
	}

	generator.Prepare(cfg)
	generator.GenerateGinCode(analyzer)
	generator.GenerateTypeScriptDTOs(analyzer)

	return nil
}
