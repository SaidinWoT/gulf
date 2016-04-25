package stream

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

const b string = "text"

var names = []string{"Test", "Pizza"}

func testSrcs() []ReadNamer {
	rns := make([]ReadNamer, len(names))
	for i, name := range names {
		rns[i] = NamedReader{
			Reader:     bytes.NewBufferString(name),
			NameString: name,
		}
	}
	return rns
}

func Prepend(b string) Transform {
	return func(s Stream) Stream {
		t := make(chan ReadNamer, len(s))
		go func() {
			for rn := range s {
				t <- NamedReader{
					Reader:     io.MultiReader(bytes.NewBufferString(b), rn),
					NameString: rn.Name(),
				}
			}
			close(t)
		}()
		return t
	}
}

func TestSrc(t *testing.T) {
	rns := testSrcs()
	s := Src(rns...)

	for rn := range s {
		if !in(names, rn.Name()) {
			t.Error("Member not given correct name.")
		}
		msg, err := ioutil.ReadAll(rn)
		if err != nil {
			t.Error(err)
		}
		if rn.Name() != string(msg) {
			t.Error("Reader not properly instantiated.")
		}
	}
}

func TestPipe(t *testing.T) {
	rns := testSrcs()
	s := Src(rns...).Pipe(Prepend(b))
	for rn := range s {
		cmp := b + rn.Name()
		msg, err := ioutil.ReadAll(rn)
		if err != nil {
			t.Error(err)
		}
		if cmp != string(msg) {
			t.Errorf("Transformation did not work. %s != %s", cmp, string(msg))
		}
	}
}

func TestDest(t *testing.T) {
	ws := make(map[string]io.Writer, len(names))
	for _, name := range names {
		ws[name] = new(bytes.Buffer)
	}

	rns := testSrcs()
	err := Src(rns...).Dest(func(name string) (io.Writer, error) {
		return ws[name], nil
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

func in(ss []string, p string) bool {
	for _, s := range ss {
		if s == p {
			return true
		}
	}
	return false
}
