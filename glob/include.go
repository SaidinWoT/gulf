package glob

import (
	"log"
	"path/filepath"
)

// include constructs a map of files to be included by running the set of patterns through the provided globFn.
// Filepaths are made absolute and cleaned to eliminate most possible duplications.
func include(globFn func(string) ([]string, error), patterns ...string) (includeMap, error) {
	m := make(map[string]bool)
	for _, pattern := range patterns {
		include := true
		if pattern[0] == '!' {
			pattern = pattern[1:]
			include = false
		}
		p, err := globFn(pattern)
		if err != nil {
			return m, err
		}
		for _, s := range p {
			s, err := filepath.Abs(s)
			if err != nil {
				log.Println(err)
			}
			s = filepath.Clean(s)
			m[s] = include
		}
	}
	return m, nil
}

func (m includeMap) list() []string {
	var fs []string
	for s, include := range m {
		if include {
			fs = append(fs, s)
		}
	}
	return fs
}
