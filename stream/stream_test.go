package stream_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"testing"

	. "github.com/SaidinWoT/gulf/stream"
)

const b = "Text"

var names = []string{"Test", "Pizza"}

func SrcMap() map[string]io.Reader {
	m := make(map[string]io.Reader)
	for _, name := range names {
		m[name] = bytes.NewBufferString(name)
	}
	return m
}

func Prepend(b string) Transform {
	return func(s Stream) Stream {
		var t Stream
		for _, m := range s.M {
			m.Reader = io.MultiReader(bytes.NewBufferString(b), m.Reader)
			t.M = append(t.M, m)
		}
		return t
	}
}

func TestSrc(t *testing.T) {
	m := SrcMap()
	s := Src(m)
	if len(s.M) != 2 {
		t.Error("Improper number of Members created.")
	}
	for _, m := range s.M {
		if m.Name != names[0] && m.Name != names[1] {
			t.Error("Member not given correct name.")
		}
		msg, err := ioutil.ReadAll(m.Reader)
		if err != nil {
			t.Error(err)
		}
		if m.Name != string(msg) {
			t.Error("Reader not properly instantiated.")
		}
	}
}

func TestPipe(t *testing.T) {
	m := SrcMap()
	s := Src(m).Pipe(Prepend(b))
	for _, m := range s.M {
		cmp := b + m.Name
		msg, err := ioutil.ReadAll(m.Reader)
		if err != nil {
			t.Error(err)
		}
		if cmp != string(msg) {
			t.Errorf("Transformation did not work. %s != %s", cmp, string(msg))
		}
	}
	m = SrcMap()
	s = Src(m)
	s.Err = errors.New("")
	s = s.Pipe(Prepend(b))
	for _, m := range s.M {
		msg, err := ioutil.ReadAll(m.Reader)
		if err != nil {
			t.Error(err)
		}
		str := string(msg)
		cmp := b + m.Name
		if cmp == str {
			t.Error("Transformation not aborted by Stream error.")
		}
		if m.Name != str {
			t.Errorf("Reader not preserved during Stream error. %s != %s", b, str)
		}
	}
}

func TestDest(t *testing.T) {
	ws := make(map[string]io.Writer, len(names))
	for _, n := range names {
		ws[n] = new(bytes.Buffer)
	}
	m := SrcMap()
	err := Src(m).Dest(func(name string) io.Writer {
		return ws[name]
	})
	if err != nil {
		t.Error(err)
	}
	for n, w := range ws {
		r := bytes.NewBuffer(w.(*bytes.Buffer).Bytes())
		msg, err := ioutil.ReadAll(r)
		if err != nil {
			t.Error(err)
		}
		if n != string(msg) {
			t.Errorf("Reader not appropriately written to Writer. %s != %s", n, string(msg))
		}
	}
}
