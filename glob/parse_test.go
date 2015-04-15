package glob_test

import (
	"testing"

	. "github.com/SaidinWoT/gulf/glob"
)

type ancestors struct {
	patterns []string
	ancestor string
}

var Ancestors = []ancestors{
	{[]string{}, ""},
	{[]string{"foo/bar/baz"}, "foo/bar/"},
	{[]string{"foo/bar/baz", "foo/bar", "foo/bar/qux"}, "foo/"},
	{[]string{"/home/glob/foo", "/home/glob/bar", "/home/glob/baz"}, "/home/glob/"},
	{[]string{"/", "/home", "/home/glob", "/home/glob/foo"}, "/"},
	{[]string{"/home/glob/foo", "/home/glob", "/home", "/"}, "/"},
	{[]string{"", "home", "home/glob", "home/glob/foo"}, ""},
	{[]string{"home/glob/foo", "home/glob", "home", ""}, ""},
	{[]string{"foo/bar", "foo/bar", "foo/bar"}, "foo/"},
	{[]string{"foo/", "foo/bar", "foo/bar/baz"}, "foo/"},
}

func TestAncestor(t *testing.T) {
	for _, test := range Ancestors {
		a := Ancestor(test.patterns...)
		if a != test.ancestor {
			t.Errorf(`Ancestor: "%s" != "%s"`, a, test.ancestor)
		}
	}
}

type prefixes struct {
	pattern, exclude, prefix string
}

var Prefixes = []prefixes{
	{"/abc", "", "/"},
	{"/abc/", "", "/abc/"},
	{"/abc/def/ghi", "beh", "/"},
	{"/abc/def/ghi", "jkl", "/abc/def/"},
	{"/abc/def/ghi", "dxz", "/abc/"},
	{"/abc/def/ghi", "zxa", "/"},
}

func TestPrefix(t *testing.T) {
	for _, test := range Prefixes {
		p := Prefix(test.pattern, []rune(test.exclude))
		if p != test.prefix {
			t.Errorf(`Prefix: "%s" != "%s"`, p, test.prefix)
		}
	}
}
