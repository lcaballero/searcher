package index

import (
	"os"
	"sync"
	"time"
)

type Info struct {
	Name     string
	Size     int64
	Mode     os.FileMode
	ModeTime time.Time
	IsDir    bool
}

func NewInfo(info os.FileInfo) *Info {
	return &Info{
		Name:     info.Name(),
		Size:     info.Size(),
		Mode:     info.Mode(),
		ModeTime: info.ModTime(),
		IsDir:    info.IsDir(),
	}
}

type Dex struct {
	Path  string
	Bytes []byte
	Info  *Info
}

func NewDex(path string, bytes []byte, info os.FileInfo) *Dex {
	return &Dex{
		Path:  path,
		Bytes: bytes,
		Info:  NewInfo(info),
	}
}

func (d *Dex) Len() int {
	return len(d.Bytes)
}

type Indexable struct {
	*Dex
	wait  *sync.WaitGroup
	Count int
}

func NewIndexable(dex *Dex, wg *sync.WaitGroup) *Indexable {
	return &Indexable{
		Dex:  dex,
		wait: wg,
	}
}

func (n Indexable) Done() {
	n.wait.Done()
}
