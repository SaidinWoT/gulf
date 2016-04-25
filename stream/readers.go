package stream

import "io"

// NamedReader is a simple implementation of ReadNamer, holding a
// NameString alongside an io.Reader.
type NamedReader struct {
	io.Reader
	NameString string
}

func (nr NamedReader) Name() string {
	return nr.NameString
}

type readEOFCloser struct {
	io.Reader
}

// ReadEOFCloser returns r wrapped such that it will automatically
// attempt to close r when it returns io.EOF.
func ReadEOFCloser(r io.Reader) io.Reader {
	return readEOFCloser{
		Reader: r,
	}
}

func (r readEOFCloser) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if err == io.EOF {
		if c, ok := r.Reader.(io.Closer); ok {
			c.Close()
		}
	}
	return n, err
}

type teeReadCloser struct {
	io.Reader
	io.Closer
}

// TeeReadCloser returns a io.ReadCloser which writes every read from r to wc.
// The returned io.ReadCloser calls wc's Close.
func TeeReadCloser(r io.Reader, wc io.WriteCloser) io.ReadCloser {
	return teeReadCloser{
		Reader: io.TeeReader(r, wc),
		Closer: wc,
	}
}
