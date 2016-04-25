package stream

import "io"

// Transform is any function that modifies a Stream.
type Transform func(Stream) Stream

// Pipe applies the provided Transforms to the current state of the Stream.
func (s Stream) Pipe(ts ...Transform) Stream {
	for _, t := range ts {
		s = t(s)
	}
	return s
}

// Fork returns two functionally identical Streams.
//
// Calls to Dest and Wait for the returned Streams must run in parallel.
func (s Stream) Fork() (Stream, Stream) {
	c, d := make(chan ReadNamer, len(s)), make(chan ReadNamer, len(s))
	go func() {
		for r := range s {
			pr, pw := io.Pipe()
			c <- NamedReader{
				Reader:     ReadEOFCloser(TeeReadCloser(r, pw)),
				NameString: r.Name(),
			}
			d <- NamedReader{
				Reader:     pr,
				NameString: r.Name(),
			}
		}
		close(c)
		close(d)
	}()
	return c, d
}

// Split divides s into two streams based on a predicate function p.
// The first Stream contains all ReadNamers for which p returns true,
// the second contains the remaining ReadNamers.
func (s Stream) Split(p func(ReadNamer) bool) (Stream, Stream) {
	c, d := make(chan ReadNamer, len(s)), make(chan ReadNamer, len(s))
	go func() {
		for r := range s {
			if p(r) {
				c <- r
			} else {
				d <- r
			}
		}
		close(c)
		close(d)
	}()
	return c, d
}

// Merge returns a Transform that combines two Streams.
func Merge(s Stream) Transform {
	return func(t Stream) Stream {
		c := make(chan ReadNamer, len(s)+len(t))
		go func() {
			for {
				select {
				case r, ok := <-s:
					if !ok {
						s = nil
					} else {
						c <- r
					}
				case r, ok := <-t:
					if !ok {
						t = nil
					} else {
						c <- r
					}
				}
				if s == nil && t == nil {
					break
				}
			}
			close(c)
		}()
		return c
	}
}
