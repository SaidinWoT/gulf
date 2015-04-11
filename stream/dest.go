package stream

import (
	"io"
	"io/ioutil"
	"sync"
)

// The ErrIntentional type serves as a specific error type to be set when a Member
// should be removed from further processing and not be written to any provided destinations.
// It serves as a marker that there was not a critical error ending the Member's usage.
type ErrIntentional struct {
	error
}

// IntentionalErr simply wraps up an error as an ErrIntentional,
// allowing use and preservation of custom errors as needed.
func IntentionalErr(e error) error {
	return ErrIntentional{e}
}

// Dest writes the results of the Stream to io.Writers returned by supplying
// fn with the names of the Members making up the Stream.
//
// It returns the first error it encounters on either the Stream or its Members,
// with the Stream's error taking precedence.
//
// Any Member with an Err set will not be written.
// If the Err in the overall Stream is set, none of its Members will be written.
// In either case, any non-nil Readers will be fully drained before program termination.
func (s Stream) Dest(fn func(string) io.Writer) error {
	return s.waitOn(func(m Member) {
		w := fn(m.Name)
		if w != nil {
			io.Copy(w, m.Reader)
		}
	})
}

// Wait waits for the Members to go through their series of transformations.
//
// It returns the first error it encounters on either the Stream or its Members,
// with the Stream's error taking precedence.
//
// All output from the final Pipe is discarded; use Dest to write it somewhere.
func (s Stream) Wait() error {
	return s.waitOn(drain)
}

func drain(m Member) {
	io.Copy(ioutil.Discard, m.Reader)
}

func (s Stream) waitOn(f func(Member)) error {
	if s.Err != nil {
		f = drain
	}
	w := new(sync.WaitGroup)
	w.Add(len(s.M))
	for _, m := range s.M {
		if _, intentional := m.Err.(ErrIntentional); s.Err == nil && m.Err != nil && !intentional {
			s.Err = m.Err
		}
		if m.Reader == nil {
			continue
		}
		fn := f
		if m.Err != nil {
			fn = drain
		}
		go func(m Member) {
			fn(m)
			w.Done()
		}(m)
	}
	w.Wait()
	return s.Err
}
