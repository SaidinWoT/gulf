package util

import (
	"io"
	"os"
	"path/filepath"

	"github.com/SaidinWoT/gulf/glob"
)

// SrcFiles opens the provided filenames and returns a map of names to files.
// The names are relative to the closest ancestor path element.
func SrcFiles(filenames ...string) map[string]io.Reader {
	base := glob.Ancestor(filenames...)
	return SrcFilesAt(base, filenames...)
}

// SrcFilesAt opens the provided filenames and returns a map of names to files.
// The names are relative to the provided base.
// If the filename cannot be made relative to the base, the full filename is used.
func SrcFilesAt(base string, filenames ...string) map[string]io.Reader {
	m := make(map[string]io.Reader)
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			continue
		}
		fi, err := f.Stat()
		if err != nil || fi.IsDir() {
			f.Close()
			continue
		}
		path, err := filepath.Rel(base, filename)
		if err != nil {
			path = filename
		}
		m[path] = f
	}
	return m
}

// At returns a function which creates and returns a new file in dir.
// The file's name is set to the string argument to the function.
func At(dir string) func(string) io.Writer {
	return func(name string) io.Writer {
		f, err := os.Create(filepath.Join(dir, name))
		if err != nil {
			f = nil
		}
		return f
	}
}
