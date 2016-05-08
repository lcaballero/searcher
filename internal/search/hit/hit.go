package hit

import (
	"bytes"
	"fmt"
	"strconv"
)

// A Hit represents a location in a byte slice, it is made of the byte slice
// the start and end offsets of the location.
type Hit struct {
	bytes []byte
	start int
	end   int
}

// Hits turns a list of bounds into an Hit array each with the given backing
// array of bytes.
func Hits(filename string, bytes []byte, bounds [][]int) (*FileHits, error) {
	hits := NewFileHits(filename)
	var firstErr error = nil
	for _, locs := range bounds {
		h, err := NewHit(bytes, locs)
		if err != nil && firstErr == nil {
			firstErr = err
			continue
		}
		hits.Add(h)
	}
	return hits, firstErr
}

// NewHit creates a hit based on the given byte slice, and the location
// bounds in a two element int slice.
func NewHit(bytes []byte, bounds []int) (*Hit, error) {
	if bytes == nil || len(bytes) <= 0 {
		return nil, fmt.Errorf("The backing bytes cannot be nil or empty.")
	}
	if len(bounds) != 2 {
		return nil, fmt.Errorf("Cannot create a hit without start and end.")
	}
	var start, end int = bounds[0], bounds[1]

	if start > end {
		return nil, fmt.Errorf("Bounds are not in order")
	}
	if start < 0 || end < 0 || start >= len(bytes) || end >= len(bytes) {
		return nil, fmt.Errorf("Bounds are not withing byte offsets (%d, %d)", start, end)
	}
	h := &Hit{
		bytes: bytes,
		start: start,
		end:   end,
	}
	return h, nil
}

// String provides a simplistic string representation of the hit.
func (h *Hit) String() string {
	return fmt.Sprintf("[%d,%d] %s", h.start, h.end, string(h.bytes[h.start:h.end]))
}

// WriteToBuffer outputs this Hit to the buffer.
func (h *Hit) WriteToBuffer(buf *bytes.Buffer) {
	buf.WriteRune('(')
	buf.WriteString(strconv.Itoa(h.start))
	buf.WriteRune(',')
	buf.WriteString(strconv.Itoa(h.end))
	buf.WriteRune(')')
	buf.WriteRune(' ')
	buf.Write(h.bytes[h.start:h.end])
}
