package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Description:          "Nornir allows codegeneration through comment-based annotations for GIN",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			InitCommand(),
			GenerateCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
