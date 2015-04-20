package task

import (
	"errors"
	"strings"
	"sync"
)

// ErrNoName is returned if a task has no name after flags have been removed.
var ErrNoName = errors.New("You defined a task with an empty name (after flag removal)")

// ErrSameName indicates that a task already exists with the provided name.
type ErrSameName struct {
	name string
}

func (e ErrSameName) Error() string {
	return "You defined multiple tasks named: " + e.name
}

// ErrCycle indicates that a dependency cycle would be created by registering the provided task.
type ErrCycle []string

func (e ErrCycle) Error() string {
	return "Cycle identified: " + strings.Join(e, ", ")
}

// ErrTaskNotExist is returned to indicate that a requested name was not found in the task map.
type ErrTaskNotExist struct {
	name string
}

func (e ErrTaskNotExist) Error() string {
	return "Task " + e.name + " does not exist."
}

// ErrExec indicates any failures encountered while executing a task.
type ErrExec struct {
	sync.Mutex
	Task error            // The error returned by the task itself.
	Req  map[string]error // Errors returned by required dependencies.
	Opt  map[string]error // Errors returned by optional dependencies. For introspection only.
}

func newErrExec() *ErrExec {
	return &ErrExec{
		Req: make(map[string]error),
		Opt: make(map[string]error),
	}
}

// Add adds an error for the named dependency.
func (e *ErrExec) Add(name string, err error, optional bool) {
	e.Lock()
	m := e.Req
	if optional {
		m = e.Opt
	}
	m[name] = err
	e.Unlock()
}

// Failed reports whether the task has had any of its required dependencies fail.
func (e *ErrExec) Failed() bool {
	return len(e.Req) > 0
}

func (e *ErrExec) Error() string {
	if e.Task != nil {
		return e.Task.Error()
	}
	var failed []string
	for name := range e.Req {
		failed = append(failed, name)
	}
	return "Failed Dependencies: " + strings.Join(failed, ", ")
}
