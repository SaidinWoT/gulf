package stream

import (
	"io"
	"io/ioutil"
	"sync"
)

// Dest writes the results of the Stream to io.Writers returned by supplying
// fn with the names of the ReadNamers making up the Stream.
//
// It returns any error encountered in creating the io.Writer or writing to it.
func (s Stream) Dest(fn func(string) (io.Writer, error)) (err error) {
	e := make(chan error)
	wg := new(sync.WaitGroup)
	defer func() {
		wg.Wait()
		select {
		case err = <-e:
		default:
		}
		close(e)
	}()
	for {
		select {
		case r, ok := <-s:
			if !ok {
				return nil
			}
			w, err := fn(r.Name())
			if err != nil {
				return err
			}
			wg.Add(1)
			go func(w io.Writer, r io.Reader) {
				_, err := io.Copy(w, r)
				if err != nil {
					e <- err
				}
				wg.Done()
			}(w, r)
		case err := <-e:
			return err
		}
	}
}

// Wait waits for all ReadNamers to go through their series of transformations.
//
// All output from the final Pipe is discarded; use Dest to write it somewhere.
func (s Stream) Wait() error {
	return s.Dest(func(_ string) (io.Writer, error) {
		return ioutil.Discard, nil
	})
}
