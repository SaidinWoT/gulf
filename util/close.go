package util

import (
	"errors"
	"io"

	"github.com/SaidinWoT/gulf/stream"
)

// IntentionallyClosed is an ErrIntentional error.
// Its presence indicates that a Member's closure should not be considered an error.
var IntentionallyClosed = stream.IntentionalErr(errors.New("Member was intentionally closed."))

type nilReader struct {
	r io.Reader
}

// NilReader provides a type which discards everything read until an error is reported.
// It returns 0, nil until an error, at which point it will return 0, io.EOF.
//
// NilReader should only be used when fully discarding a Member.
func NilReader(r io.Reader) io.Reader {
	return nilReader{r: r}
}

func (r nilReader) Read(p []byte) (int, error) {
	b := make([]byte, len(p))
	_, err := r.r.Read(b)
	if err != nil {
		err = io.EOF
	}
	return 0, err
}

func (r nilReader) Close() error {
	if c, ok := r.r.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// Close marks a Member as intentionally closed and replaces its Reader with a NilReader.
func Close(m stream.Member) stream.Member {
	return CloseSaveError(m, IntentionallyClosed)
}

// CloseSaveError converts a Member to a NilReader while preserving the provided error e for future use.
// Note that the error is not wrapped as an ErrIntentional.
func CloseSaveError(m stream.Member, e error) stream.Member {
	return stream.Member{
		Name:   m.Name,
		Reader: NilReader(m.Reader),
		Err:    e,
	}
}
