package stream

import "io"

// Transform is any function that takes in a Stream, transforms it, and outputs a Stream.
//
// Any modification of the io.Readers in the Members of the Stream should make use of
// one of the Mutate functions, or be modeled after them.
type Transform func(Stream) Stream

// Pipe applies the provided Transform to the current state of the stream.
// It short-circuits if the stream has an error.
func (s Stream) Pipe(f Transform) Stream {
	if s.Err != nil {
		return s
	}
	return f(s)
}

// Fork returns two functionally identical Streams.
//
// Safe usage of the Streams requires that calls to Wait-like functions run in parallel.
func (s Stream) Fork() (Stream, Stream) {
	var t Stream
	tm := make([]Member, len(s.M))
	for i, m := range s.M {
		r, w := io.Pipe()
		s.M[i] = Member{
			Name:   m.Name,
			Reader: EOFCloseReader(ReadTeeCloser(m.Reader, w)),
		}
		tm[i] = Member{
			Name:   m.Name,
			Reader: r,
		}
	}
	t.M = tm
	t.Err = s.Err
	return s, t
}

// Split divides s into two Streams.
// The first Stream will contain all the Members for which f returns true, the second with the remaining Members.
func (s Stream) Split(f func(m Member) bool) (Stream, Stream) {
	var t, u Stream
	for _, m := range s.M {
		if f(m) {
			t.M = append(t.M, m)
		} else {
			u.M = append(u.M, m)
		}
	}
	t.Err, u.Err = s.Err, s.Err
	return t, u
}

// Merge returns a function that combines two Streams, appending t's Members to s.
//
// Merge is primarily a convenience function, provided to have a definitive way to rejoin streams created by Fork (to safely Wait) or Split.
func Merge(s Stream) Transform {
	return func(t Stream) Stream {
		s.M = append(s.M, t.M...)
		if s.Err == nil {
			s.Err = t.Err
		}
		return s
	}
}

type eofCloseReader struct {
	r io.Reader
}

func (r eofCloseReader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	if err == io.EOF {
		if c, ok := r.r.(io.Closer); ok {
			c.Close()
		}
	}
	return n, err
}

// EOFCloseReader creates a Reader which tries to close itself on EOF.
func EOFCloseReader(r io.Reader) io.Reader {
	return eofCloseReader{
		r: r,
	}
}

type readTeeCloser struct {
	io.Reader
	io.Closer
}

// ReadTeeCloser creates a TeeReader which writes every read from r to w.
// The returend ReadTeeCloser's Close method calls w's Close.
//
// ReadTeeCloser allows a TeeReader to write to an io.PipeWriter and close it afterwards without having to separately maintain a reference.
func ReadTeeCloser(r io.Reader, w io.WriteCloser) io.ReadCloser {
	return readTeeCloser{
		io.TeeReader(r, w),
		w,
	}
}
