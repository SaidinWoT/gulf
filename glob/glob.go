// Package glob provides an extended globbing system.
// It suplements path/filepath's globs with globstar, pattern negation, and dotfile filtering.
package glob

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type includeMap map[string]bool

// SpecialRunes identifies the runes that have special meaning in patterns accepted by Glob.
var SpecialRunes = []rune("*?")

// Match returns true if name matches the provided pattern.
// The pattern syntax is identical to that of path.Match, except as follows:
//	'*' does not match dotfiles unless explicitly preceded by a dot.
//	'**/' (globstar) matches 0 or more directories (not including dot directories).
// globstar only has a special meaning if it is the only pattern in its containing element
// (that is, it is either preceded by the beginning of the string or a filepath separator).
// Otherwise, it is treated simply as two sequential globs (which will then be condensed to a single glob).
func Match(pattern, name string) (bool, error) {
	re, err := makeRegexp(pattern)
	if err != nil {
		return false, err
	}
	abs, err := filepath.Abs(name)
	if err != nil {
		return false, err
	}
	return re.MatchString(filepath.ToSlash(abs)), nil
}

// Glob returns the list of filenames which match the provided pattern.
// The syntax is the same as for Match.
func Glob(pattern string) ([]string, error) {
	// Lacking a globstar, use filepath.Glob with tedious dotfile filtering
	if strings.Index(pattern, "**") == -1 {
		patterns := negateDotfiles(pattern)
		m, err := include(filepath.Glob, patterns...)
		return m.list(), err
	}
	return globstar(pattern)
}

// negateDotfiles appends a filter for each glob in pattern that does not explicitly follow a dot character.
// Each filter will be identical to pattern except for an explicit dot before one of the globs.
func negateDotfiles(pattern string) []string {
	patterns := []string{pattern}
	pieces := strings.Split(pattern, "*")
	l := len(pieces)
	for i := 0; i < l-1; i++ {
		p := pieces[i]
		if n := len(p); n != 0 && p[n-1] == '.' {
			continue
		}
		ps := make([]string, l)
		copy(ps, pieces)
		ps[i] += "."
		ps[0] = "!" + ps[0]
		patterns = append(patterns, strings.Join(ps, "*"))
	}
	return patterns
}

// globstar does a filesystem walk rooted at the longest definite prefix of pattern.
// It returns the list of filenames that match the pattern as dictated by the syntax of Match.
func globstar(pattern string) ([]string, error) {
	var s []string
	re, err := makeRegexp(pattern)
	if err != nil {
		return s, err
	}
	filepath.Walk(Prefix(pattern, SpecialRunes), func(path string, info os.FileInfo, _ error) error {
		if info == nil {
			log.Println("Malformed directory", path)
			return filepath.SkipDir
		}
		path, err = filepath.Abs(path)
		if err != nil {
			return nil
		}
		if re.MatchString(filepath.ToSlash(path)) {
			s = append(s, path)
		}
		return nil
	})
	return s, nil
}

// makeRegexp compiles a regular expression that will match a filepath containing a globstar.
// The resulting regex expects matches to have passed through filepath.ToSlash.
func makeRegexp(p string) (*regexp.Regexp, error) {
	p, err := filepath.Abs(p)
	if err != nil {
		return nil, err
	}
	p = filepath.ToSlash(p)
	p = regexify(p)
	p = "^" + p + "$"
	return regexp.Compile(p)
}

// regexify converts a pattern string fit for filepath.Match into a pattern for regexp.MatchString.
// The resulting regexp should have the same results as filepath.Match
// (with the differences mentioned in the documentation for Match).
func regexify(p string) string {
	// The worst case pattern is a cycle of globs followed by a file separator.
	// This will replace every 2 characters with 16 characters, so 8*len(p) will suffice.
	s := make([]byte, 8*len(p))
	j := 0
	onseparator := true
	inrange := false
	l := len(p)
	for i := 0; i < l; i++ {
		for inrange && i < l {
			switch p[i] {
			case '\\':
				s[j] = p[i]
				i++
				j++
			case ']':
				inrange = false
			}
			s[j] = p[i]
			i++
			j++
		}
		keepchar := true
		switch p[i] {
		case '\\':
			s[j] = p[i]
			i++
			j++
		case '.':
			s[j] = '\\'
			j++
			if l > i+1 && p[i+1] == '*' {
				s[j] = '.'
				i++
				j++
				goto Suffix
			}
		case '*':
			if !onseparator {
				goto Suffix
			}
			if l > i+3 && p[i:i+3] == "**/" {
				_ = append(s[0:j], []byte("(?:[^./][^/]*/)*")...)
				i += 2
				// Condense sequential globstars
				for l > i+3 && p[i+1:i+4] == "**/" {
					i += 3
				}
				j += 16
			} else {
				_ = append(s[0:j], []byte("(?:[^./][^/]*)?")...)
				// Condense sequential globs
				// (?s cannot easily combine with dotfile-excluding globs)
				for l > i+1 && p[i+1] == '*' {
					i++
				}
				j += 15
			}
			keepchar = false
		case '?':
			goto Suffix
		case '[':
			inrange = true
		case '+', '(', ')', '|', '{', '}', '^', '$':
			s[j] = '\\'
			j++
		}
		if keepchar {
			s[j] = p[i]
			j++
		}
	Separator:
		onseparator = p[i] == '/'
		continue
	Suffix:
		// Condense sequential *s and ?s into a single regex
		n, suffix := wildcardSuffix(p[i:])
		_ = append(s[0:j], []byte("[^/]"+suffix)...)
		i += n - 1
		j += 4 + len(suffix)
		goto Separator
	}
	return string(s[0:j])
}

// wildcardSuffix returns the number of sequential wildcards and an appropriate regexp suffix for their characteristics.
func wildcardSuffix(p string) (skip int, suffix string) {
	var req int
	var star bool
	for _, b := range p {
		switch b {
		case '?':
			req++
			skip++
		case '*':
			star = true
			skip++
		default:
			break
		}
	}
	if req > 1 {
		var comma string
		if star {
			comma = ","
		}
		suffix = "{" + strconv.Itoa(req) + comma + "}"
	} else if req == 1 && star {
		suffix = "+"
	} else if star {
		suffix = "*"
	}
	return
}
