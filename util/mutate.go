package util

import (
	"io"

	. "github.com/SaidinWoT/gulf/stream"
)

// Mutate wraps f, returning a Transform that implements the main streaming functionality.
//
// The Reader from each Member of the input Stream is copied to a PipeWriter; the associated PipeReader is passed to f.
// The Reader returned by f must make use of the PipeReader - failing to do so destroys the input and may result in deadlock.
//
// The error returned by f will be set to the Member being mutated.
// Such errors generally stop any further transformation and should thus be returned only when necessary.
func Mutate(f func(io.Reader) (io.Reader, error)) Transform {
	return mutate(
		func(m Member, r io.Reader) Member {
			m.Reader, m.Err = f(r)
			return m
		})
}

// MutateMembers serves the same purpose as Mutate but allows the function to modify an entire Member.
// The Reader in the provided member is still the PipeReader and not the actual Reader from the input
//
// Generally the name should be left unchanged by f - its inclusion is to allow modifications of the Reader based on name.
// However, functionality that has good reason to change the name is not prohibited from doing so.
//
// Errors that cannot be recovered from should be set directly in the returned Member.
// Note that this will generally stop any further transformation of that Member and should thus be a last resort.
func MutateMembers(f func(Member) Member) Transform {
	return mutate(
		func(m Member, r io.Reader) Member {
			m.Reader = r
			return f(m)
		})
}

func mutate(f func(Member, io.Reader) Member) Transform {
	return func(s Stream) Stream {
		var t Stream
		t.M = make([]Member, len(s.M))
		for i, m := range s.M {
			if m.Err != nil {
				t.M[i] = m
				continue
			}
			r, w := io.Pipe()
			go pipe(w, m.Reader)
			t.M[i] = f(m, r)
		}
		return t
	}
}

func pipe(dst *io.PipeWriter, src io.Reader) {
	_, err := io.Copy(dst, src)
	dst.CloseWithError(err)
	if c, ok := src.(io.Closer); ok {
		c.Close()
	}
}
