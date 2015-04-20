package watch

import (
	"time"

	"github.com/SaidinWoT/gulf/glob"
	"github.com/SaidinWoT/gulf/task"
)

// The Set type embeds a task.Set and extends it with the ability to
// execute tasks based on filesystem events.
type Set struct {
	*task.Set
	match func(string, string) (bool, error)
	glob  func(string) ([]string, error)
	wm    map[string][]string
	delay time.Duration
}

// New returns a new Set, initialized to use glob.Match and glob.Glob.
func New() *Set {
	return &Set{
		Set:   task.New(),
		match: glob.Match,
		glob:  glob.Glob,
		wm:    make(map[string][]string),
		delay: 10 * time.Millisecond,
	}
}

// The Option type is a function that modifies a Watch.
type Option func(s *Set) error

// SetOption modifies s with Options provided.
func (s *Set) SetOption(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return err
		}
	}
	return nil
}

// Globber sets the globbing function used by the Set to identify files to watch.
//
// The default is glob.Glob
func Globber(fn func(string) ([]string, error)) Option {
	return func(s *Set) error {
		s.glob = fn
		return nil
	}
}

// Matcher sets the matching function used by the Set to assert that
// events on a file should trigger tasks.
//
// The default is glob.Match
func Matcher(fn func(string, string) (bool, error)) Option {
	return func(s *Set) error {
		s.match = fn
		return nil
	}
}

// Delay sets the delay used by Watch to identify unique file events.
// All file events on a single file name (not inode) that take place within d time
// of each other will be considered a single file event by Watch.
// This is a debouncing mechanism, provided to trigger tasks only once despite
// most text editor's usage of multiple file events while writing a file.
//
// The default is 10 milliseconds
func Delay(d time.Duration) Option {
	return func(s *Set) error {
		s.delay = d
		return nil
	}
}
