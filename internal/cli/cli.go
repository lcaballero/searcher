package cli

import (
	cmd "github.com/codegangsta/cli"
	"github.com/lcaballero/searcher/internal"
)

func NewCli() *cmd.App {
	app := cmd.NewApp()
	app.Name = "Denorm"
	app.Version = "0.0.1"
	app.Usage = "The 'Denormalization' service with commands and tools for this service."
	app.Commands = []cmd.Command{
		searchCommand(),
		indexCommand(),
		promptCommand(),
	}
	app.Action = internal.Searching
	return app
}

func promptCommand() cmd.Command {
	return cmd.Command{
		Name:   "prompt",
		Usage:  "Starts the search prompt.",
		Action: internal.Prompt,
		Flags: []cmd.Flag{
			cmd.StringFlag{
				Name:  "index",
				Value: "index.bin",
				Usage: "File name of the index.",
			},
		},
	}
}

func searchCommand() cmd.Command {
	return cmd.Command{
		Name:   "re",
		Usage:  "Requests the lookup and displays it.",
		Action: internal.Searching,
		Flags: []cmd.Flag{
			cmd.StringFlag{
				Name:  "index",
				Value: "index.bin",
				Usage: "File name of the index.",
			},
		},
	}
}

func indexCommand() cmd.Command {
	return cmd.Command{
		Name:   "index",
		Usage:  "Traverses the root",
		Action: internal.Indexing,
		Flags: []cmd.Flag{
			cmd.StringFlag{
				Name:  "root",
				Value: ".files/testdata",
				Usage: "Uses the given directory as the root for indexing.",
			},
			//TODO: Not implemented, but stubbed out here to remind myself
			cmd.StringFlag{
				Name:  "extensions",
				Value: "js,txt,xml,java,sh,html,css,c,cpp,pl,xsl,properties",
				Usage: "Uses the given directory as the root for indexing.",
			},
		},
	}
}
