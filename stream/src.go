// Package stream implements a simple system to sequentially modify a set of named readers.
package stream

import "io"

// A ReadNamer is an io.Reader accompanied by a Name.
type ReadNamer interface {
	io.Reader
	Name() string
}

// A Stream is a collection of ReadNamers.
//
// Streams are intended to be sent trough a series of transformations with Pipe,
// usually to ultimately be written to files with Dest.
type Stream <-chan ReadNamer

// Src creates a Stream populated with the provided ReadNamers.
func Src(rs ...ReadNamer) Stream {
	c := make(chan ReadNamer, len(rs))
	for _, r := range rs {
		c <- r
	}
	close(c)
	return c
}
