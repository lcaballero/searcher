package hit

import "bytes"

type FileHits struct {
	Filename string
	hits []*Hit
}

func NewFileHits(filename string) *FileHits {
	return &FileHits{
		Filename: filename,
		hits: make([]*Hit, 0),
	}
}

func (f *FileHits) Add(h *Hit) {
	f.hits = append(f.hits, h)
}

func (f *FileHits) Len() int {
	return len(f.hits)
}

func (f *FileHits) String() string {
	buf := &bytes.Buffer{}
	for _, h := range f.hits {
		h.WriteToBuffer(buf)
		buf.WriteRune('\n')
	}
	return buf.String()
}
