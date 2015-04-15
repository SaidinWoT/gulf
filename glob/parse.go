package glob

import (
	"os"
	"path/filepath"
)

// Ancestor returns the closest common ancestor of the provided paths, with the trailing file separator.
// If there is no common ancestor, an empty string is returned.
//
// Ancestor does not handle any special characters - thus
//	Ancestor("foo/*/bar", "foo/*/qux")
// will return "foo/*/".
func Ancestor(patterns ...string) string {
	if l := len(patterns); l == 0 {
		return ""
	} else if l == 1 {
		p := filepath.Dir(patterns[0])
		if p == "." {
			return ""
		} else {
			return p + string(filepath.Separator)
		}
	}
	path := patterns[0]
	rest := make([][]rune, len(patterns)-1)
	for i, s := range patterns[1:] {
		rest[i] = []rune(s)
	}
	return prefixer(path, rest, func(p rune, r []rune, i int) bool {
		return len(r) <= i || r[i] != p
	})
}

// Prefix returns all elements (plus trailing file separator) of the path up to,
// but not including, the first element containing a rune in exclude.
// If there is no such prefix, it returns the empty string.
func Prefix(path string, exclude []rune) string {
	exs := make([][]rune, len(exclude))
	for i, r := range exclude {
		exs[i] = []rune{r}
	}
	return prefixer(path, exs, func(p rune, r []rune, _ int) bool {
		return p == r[0]
	})
}

func prefixer(path string, rs [][]rune, fn func(rune, []rune, int) bool) string {
	sep := 0
	for i, p := range path {
		for _, r := range rs {
			if fn(p, r, i) {
				return path[:sep]
			}
		}
		if p == filepath.Separator {
			sep = i + 1
		}
	}
	return path[:sep]
}

// Parse returns the full list of file and directory names matched by
// the provided globbing function with the list of patterns.
// A pattern with '!' as its first character is treated as a negation.
// Patterns are processed in sequential order - negations will remove
// matches from preceding patterns, but not those which follow them.
func Parse(globFn func(string) ([]string, error), patterns ...string) []string {
	wd, err := os.Getwd()
	if err != nil {
		return nil
	}
	m, err := include(globFn, patterns...)
	l := m.list()
	for i, s := range l {
		l[i], _ = filepath.Rel(wd, s)
	}
	return l
}
