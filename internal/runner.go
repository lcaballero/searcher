package internal

import (
	"fmt"
	"os"

	cmd "github.com/codegangsta/cli"
	indexing "github.com/lcaballero/searcher/internal/index"
	"github.com/lcaballero/searcher/internal/search/exec"
	"github.com/lcaballero/searcher/internal/search/hit"
	"github.schq.secious.com/Logrhythm/GoDispatch/bench"
	"github.com/lcaballero/searcher/internal/console"
	"github.com/vrecan/death"
	"syscall"
)

func runSearch(
	index *indexing.IndexFile,
	res chan *hit.FileHits,
	start chan bool,
	complete chan *indexing.IndexFile) {

	go func() {
		tc := bench.Start()
		file := 0
		hits := 0
		for {
			select {
			case <-start:
				tc = bench.Start()
				file = 0
				hits = 0
			case dex := <-complete:
				tc.Stop()
				fmt.Printf("Index name: %s\n", dex.Filename)
				fmt.Printf("Elapsed: %s\n", tc.Elapsed())
				fmt.Printf("File hits: %d\n", file)
				fmt.Printf("Hits: %d\n", hits)
			case hts := <-res:
				if hts.Len() > 0 {
					file++
					hits += hts.Len()
				}
			}
		}
	}()
}

func Searching(ctx *cmd.Context) {
	pattern := ctx.Args().Get(0)
	indexName := ctx.String("index")

	index, err := indexing.LoadIndexFile(indexName)
	if err != nil {
		panic(err)
	}

	index.Report(os.Stdout)
	res := make(chan *hit.FileHits, 1)
	complete := make(chan *indexing.IndexFile)
	start := make(chan bool)

	runSearch(index, res, start, complete)
	exec.NewFind(pattern, index).Execute(res, start, complete)
}

func Indexing(cli *cmd.Context) {
	root := cli.String("root")
	dexing := indexing.NewIndexer(root)

	indexFile := dexing.Traverse()
	err := indexFile.Write(root, "index.bin")
	if err != nil {
		fmt.Println(err)
	}

	indexFile.Report(os.Stdout)
}

func Prompt(cli *cmd.Context) {
	prompt := console.NewPrompt()
	prompt.Start()

	death.NewDeath(syscall.SIGINT, syscall.SIGTERM).WaitForDeath(prompt)
}
