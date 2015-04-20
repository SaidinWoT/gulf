package task

import "strings"

// The Set type provides a structure to register a set of tasks and execute them.
type Set struct {
	ts  map[string]task
	err *error
}

// New returns a pointer to a Set.
func New() *Set {
	return &Set{
		ts:  make(map[string]task),
		err: new(error),
	}
}

type task struct {
	name string
	deps []string
	fn   func() error
	f    flags
}

// Flags contains the set of runes that have special meaning at the end of task names.
//	'?'	denotes an optional task - any tasks depending on it will not be prevented by the task returning an error.
//	'+' denotes a task to be executed multiple times - every instance of the task in the dependency graph will execute.
//	'*' denotes a task that is both optional and to be executed multiple times.
//
// Flags may be specified either in the task name or in a dependency.
// When specified in the task name, the special meaning will impact every instance of the task in the dependency graph.
// When specified in a dependency name, the special meaning will impact only the execution triggering the dependency.
const Flags = "?+*"

type flags struct {
	multi    bool
	optional bool
}

func parseFlags(name string) (flags, string) {
	var f flags
	l := len(name)
	for {
		name = name[:l]
		l = len(name) - 1
		switch c := name[l]; c {
		case '*':
			f.multi = true
			f.optional = true
		case '+':
			f.multi = true
		case '?':
			f.optional = true
		default:
			return f, name
		}
	}
}

// Task registers a task, correlating a function with a name and an optional set of dependencies.
// Dependencies are the string names of other tasks.
//
// Registering multiple tasks with the same name will result in only the last task being registered.
// Creating a dependency cycle registers an error in the Set, which will prevent further use of the Set.
// Any such error will also be returned.
func (s *Set) Task(name string, fn func() error, deps ...string) error {
	if *s.err != nil {
		return *s.err
	}
	f, name := parseFlags(name)
	if name == "" {
		*s.err = ErrNoName
	} else if _, exists := s.ts[name]; exists {
		*s.err = ErrSameName{name: name}
	} else if cycle, err := s.cycle(name, deps...); cycle {
		*s.err = err
	}
	if *s.err != nil {
		return *s.err
	}
	s.ts[name] = task{
		name: name,
		deps: deps,
		fn:   fn,
		f:    f,
	}
	return nil
}

func (s *Set) cycle(name string, deps ...string) (bool, error) {
	return s.search(name, name, deps...)
}

func (s *Set) search(name, curr string, deps ...string) (bool, ErrCycle) {
	if len(deps) == 0 {
		return false, nil
	}
	for _, dep := range deps {
		dep = strings.TrimRight(dep, Flags)
		if dep == name {
			return true, ErrCycle{curr}
		}
		if d, ok := s.ts[dep]; ok {
			if exists, err := s.search(name, dep, d.deps...); exists {
				return exists, append([]string{curr}, err...)
			}
		}
	}
	return false, nil
}
