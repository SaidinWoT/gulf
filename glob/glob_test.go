package glob_test

import (
	"testing"

	. "github.com/SaidinWoT/gulf/glob"
)

type matches struct {
	pattern, name string
	result        bool
}

var Matches = []matches{
	{"**/j.a", "abc/def/ghi/j.a", true},
	{"**/**/**/j.a", "abc/j.a", true},
	{"foo/*.a", "foo/abc.a", true},
	{"foo/*.a", "foo/.abc.a", false},
	{"foo/.*.a", "foo/abc.a", false},
	{"foo/.*.a", "foo/.abc.a", true},
	{"foo/.*?.a", "foo/.abc.a", true},
	{"foo/.*?.a", "foo/.a", false},
	{"foo/???.a", "foo/bar.a", true},
	{"foo/???*.a", "foo/bar.a", true},
	{"foo/????.a", "foo/bar.a", false},
	{"foo/[a-c].a", "foo/a.a", true},
	{"foo/[a-c].a", "foo/d.a", false},
	{"foo/[^a-c].a", "foo/d.a", true},
	{"foo/\\[a-c].a", "foo/[a-c].a", true},
	{"foo/\\?.a", "foo/?.a", true},
	{"foo/\\?.a", "foo/!.a", false},
	{"foo/\\*.a", "foo/*.a", true},
	{"foo/\\*.a", "foo/!.a", false},
	{"foo/\\\\.a", "foo/\\.a", true},
	{"foo/\\\\.a", "foo/!.a", false},
	{"foo/*?*.a", "foo/.a", false},
	{"foo/*?*.a", "foo/a.a", true},
	{"foo/?*.a", "foo/a.a", true},
	{"foo/?*.a", "foo/.a", false},
	{"foo/?*.a", "foo/.a.a", true},
	{"foo/a*.a", "foo/a.a.a", true},
	{"foo/*.a", "foo/.a", true},
	{"foo/***.a", "foo/.a", true},
}

func TestMatch(t *testing.T) {
	for _, test := range Matches {
		b, err := Match(test.pattern, test.name)
		if err != nil {
			t.Error(err)
		}
		if b != test.result {
			t.Errorf(`Match: "%s" matches "%s" was reported as %t`, test.pattern, test.name, b)
		}
	}
}
