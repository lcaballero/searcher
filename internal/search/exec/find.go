package exec

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/lcaballero/searcher/internal/index"
	"github.com/lcaballero/searcher/internal/search/hit"
)

type Find struct {
	pattern string
	index   *index.IndexFile
}

func NewFind(re string, dex *index.IndexFile) *Find {
	return &Find{
		pattern: re,
		index:   dex,
	}
}

func (f *Find) Execute(res chan *hit.FileHits, start chan bool, complete chan *index.IndexFile) {
	start <- true
	wg := &sync.WaitGroup{}
	wg.Add(len(f.index.Indexed))

	re, err := regexp.Compile(f.pattern)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, dex := range f.index.Indexed {
		go f.Search(re, dex, res, wg)
	}

	wg.Wait()

	complete <- f.index
}

func (f *Find) Search(re *regexp.Regexp, dex *index.Indexable, res chan *hit.FileHits, wg *sync.WaitGroup) {
	idx := re.FindAllIndex(dex.Bytes, -1)
	hits, _ := hit.Hits(dex.Path, dex.Bytes, idx)
	res <- hits
	wg.Done()
}
