package index

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Indexer struct {
	root string
	exts []string
}

func NewIndexer(root string) *Indexer {
	return &Indexer{
		root: root,
		exts: strings.Split("js,txt,xml,java,sh,html,css,c,cpp,pl,xsl,properties", ","),
	}
}

func (n *Indexer) isExtension(filename string) bool {
	//TODO: Replace from config
	for _, ex := range n.exts {
		if strings.HasSuffix(filename, ex) {
			return true
		}
	}
	return false
}

func (n *Indexer) Traverse() *IndexFile {
	wg := &sync.WaitGroup{}

	collect := make(chan *Indexable)
	indexFile := NewIndexFile()
	done := indexFile.Listen(collect)

	walk := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && n.isExtension(path) {
			wg.Add(1)
			go n.collect(path, info, collect, wg)
		}
		return nil
	}

	filepath.Walk(n.root, walk)
	wg.Wait()

	wg = &sync.WaitGroup{}
	wg.Add(1)
	done <- wg
	wg.Wait()

	return indexFile
}

func (n *Indexer) collect(path string, info os.FileInfo, fn chan *Indexable, wg *sync.WaitGroup) error {
	bb, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	dex := NewDex(path, bb, info)
	fn <- NewIndexable(dex, wg)
	return nil
}
