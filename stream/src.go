// Package stream implements a simple system to sequentially modify a set of named readers.
package stream

import "io"

// A Member is simply a named Reader.
//
// A Member will generally not be further transformed after its Err value is non-nil,
// including not being written by Dest. As a result of this,
// errors should be set only when absolutely necessary.
//
// The Name should only be changed when absolutely necessary.
type Member struct {
	Name   string
	Reader io.Reader
	Err    error
}

// A Stream is a collection of Members.
//
// Streams are intended to be sent through a series of transformations with Pipe,
// usually to ultimately be written to files with Dest.
// A Stream will not be mutated further once its Err is non-nil;
// as such, errors should only be attached to the Stream when they are impossible to recover from.
type Stream struct {
	M   []Member
	Err error
}

// Src translates a map of names to io.Readers into a Stream of Members with those values.
func Src(m map[string]io.Reader) Stream {
	var s Stream
	for name, reader := range m {
		s.M = append(s.M, Member{
			Name:   name,
			Reader: reader,
		})
	}
	return s
}

// Ordered creates a Stream with the names in ns and the associated io.Readers in m.
// The Members in the returned Stream are guaranteed to be in the same order as ns.
func Ordered(m map[string]io.Reader, ns []string) Stream {
	var s Stream
	for _, n := range ns {
		s.M = append(s.M, Member{
			Name:   n,
			Reader: m[n],
		})
	}
	return s
}
