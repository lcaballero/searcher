package main

import (
	"os"

	"github.com/lcaballero/searcher/internal/cli"
)

func main() {
	cli.NewCli().Run(os.Args)
}
