package task_test

import (
	"testing"

	. "github.com/SaidinWoT/gulf/task"
)

var returnNil = func() error { return nil }

type cycle struct {
	chain   []string
	isCycle bool
}

var Cycles = []cycle{
	{[]string{"foo", "foo"}, true},
	{[]string{"foo", "bar", "baz"}, false},
	{[]string{"foo", "bar", "baz", "foo"}, true},
}

func TestCycle(t *testing.T) {
	for _, test := range Cycles {
		s := New()
		var i int
		for i = 0; i < len(test.chain)-2; i++ {
			s.Task(test.chain[i], returnNil, test.chain[i+1])
		}
		err := s.Task(test.chain[i], returnNil, test.chain[i+1])
		if _, ok := err.(ErrCycle); ok != test.isCycle {
			t.Errorf("Cycle test for chain %v reported %t.", test.chain, ok)
		}
	}
}
