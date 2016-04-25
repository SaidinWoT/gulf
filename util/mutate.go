// Package util provides tools useful in authoring plugins for the gulf ecosystem.
package util

import (
	"io"

	"github.com/SaidinWoT/gulf/stream"
)

// Mutate wraps f, returning a Transform that implements the main streaming functionality.
func Mutate(f func(io.Reader) io.Reader) stream.Transform {
	return mutate(func(r stream.ReadNamer) stream.ReadNamer {
		return stream.NamedReader{
			Reader:     f(r),
			NameString: r.Name(),
		}
	})
}

// MutateMembers serves the same purpose as Mutate but allows the function to modify an entire Member.
func MutateMembers(f func(stream.ReadNamer) stream.ReadNamer) stream.Transform {
	return mutate(func(r stream.ReadNamer) stream.ReadNamer {
		return f(r)
	})
}

func mutate(f func(stream.ReadNamer) stream.ReadNamer) stream.Transform {
	return func(s stream.Stream) stream.Stream {
		t := make(chan stream.ReadNamer, len(s))
		go func() {
			for r := range s {
				t <- f(r)
			}
			close(t)
		}()
		return t
	}
}
