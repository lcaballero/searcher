package index

import (
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"time"

	"io"

	"github.com/fatih/color"
	"github.com/lcaballero/time-capture/bench"
)

// An IndexFile represents the set of files that were indexed, and some meta
// information about when the IndexFile was created.  Additionally, an IndexFile
// can be saved to disk, and loaded from disk.  The set of Indexable(s) is the
// data-store which searching traverses and from which 'hits' are produced.
//
// Note: and IndexFile does not update when the underlying files in the system
// are changed on the file system.  The IndexFile is a snapshot created during
// indexing.
type IndexFile struct {
	Indexed   []*Indexable
	Timestamp time.Time
	Elapsed   time.Duration
	Filename  string
	Root      string
	FileBytes int64
	capture   *bench.TimeCapture
}

// NewIndexFile creates an empty Indexfile that can be used to collect
// Indexable files through the Listen method, and then saved to disk
// with the Write method.
func NewIndexFile() *IndexFile {
	return &IndexFile{
		Indexed: make([]*Indexable, 0, 1000),
		capture: bench.Start(),
	}
}

// LoadIndexFile will decode the file and provide access to the set of
// Indexed files held in memory.
func LoadIndexFile(filename string) (*IndexFile, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	x := &IndexFile{}
	err = gob.NewDecoder(f).Decode(x)
	if err != nil {
		return nil, err
	}
	return x, nil
}

// Listen starts a Go routing to collect Indexable files which
// will make up the IndexFile prior to encoding to disk.
func (f *IndexFile) Listen(collect chan *Indexable) chan *sync.WaitGroup {
	done := make(chan *sync.WaitGroup)
	filebytes := int64(0)
	f.Timestamp = time.Now()

	go func() {
		count := 0
		for {
			select {
			case cleaner := <-done:
				f.FileBytes = filebytes
				cleaner.Done()
				return
			case dx := <-collect:
				dx.Count = count
				f.Indexed = append(f.Indexed, dx)
				count++
				filebytes += int64(dx.Len())
				dx.Done()
			}
		}
	}()

	return done
}

// Write outputs the IndexFile to disk to the given filename.  The root is
// provided for book keeping.
func (x *IndexFile) Write(root, filename string) error {
	x.Timestamp = time.Now()
	x.capture.Stop()
	x.Elapsed = x.capture.Elapsed()
	x.Filename = filename
	x.Root = root

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	return gob.NewEncoder(f).Encode(x)
}

// Report outputs some stats and information about the IndexFile, namely those
// stats that were created during the creation of the IndexFile.
func (n *IndexFile) Report(w io.Writer) {
	gr := color.New(color.FgGreen).SprintfFunc()
	fmt.Fprintln(w, "Index root at:", gr("%s", n.Root))
	fmt.Fprintln(w, "Number of files indexed:", gr("%d", len(n.Indexed)))
	fmt.Fprintln(w, "Elapsed milliseconds to create:", gr("%d", n.Elapsed/time.Millisecond))
	fmt.Fprintln(w, "Index created at:", gr("%s", n.Timestamp))
	fmt.Fprintln(w, "Index file created:", gr("%s", n.Filename))
	fmt.Fprintln(w, "Bytes indexed:", gr("%d(bytes) %d(mb)", n.FileBytes, n.FileBytes/int64(1024*1024)))
}
