// +build gulf

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/SaidinWoT/gulf/glob"
	"github.com/SaidinWoT/gulf/stream"
	"github.com/SaidinWoT/gulf/task/watch"
	"github.com/SaidinWoT/gulf/util"
)

// Gulf is a simple struct to bring together gulf's functionality.
type Gulf struct {
	s     *watch.Set
	glob  func(string) ([]string, error)
	watch bool
}

// New creates a Gulf with an empty task Set.
func New() *Gulf {
	return &Gulf{
		s:    watch.New(),
		glob: glob.Glob,
	}
}

type Option func(*Gulf) error

// SetOption modifies f with the provided Options.
func (g *Gulf) SetOption(opts ...Option) error {
	var err error
	for _, opt := range opts {
		if err = opt(g); err != nil {
			break
		}
	}
	return err
}

// Globber returns an Option that sets a Gulf's globbing function to fn.
func Globber(fn func(string) ([]string, error)) Option {
	return func(g *Gulf) error {
		return g.s.SetOption(watch.Globber(fn))
	}
}

// Matcher returns an Option that sets a Gulf's matcher function fn.
func Matcher(fn func(string, string) (bool, error)) Option {
	return func(g *Gulf) error {
		return g.s.SetOption(watch.Matcher(fn))
	}
}

// WatchDelay returns an Option that sets the delay used to detect unique file events.
func Delay(d time.Duration) Option {
	return func(g *Gulf) error {
		return g.s.SetOption(watch.Delay(d))
	}
}

// Src provides a simple wrapper around stream.Src.
// It parses patterns with f's globbing function and provides those to stream.Src.
//
// By default, a Gulf uses gulf/glob's Glob function.
// The behavior differs from filepath.Glob as follows:
// Globstar (**) matches 0 or more directories.
// Globs (*) only include dotfiles if there is an explicit dot before the glob character.
func (g *Gulf) Src(patterns ...string) stream.Stream {
	filenames := glob.Parse(g.glob, patterns...)
	m := util.SrcFiles(filenames...)
	return stream.Src(m)
}

var At = util.At

// Task adds a task to f's task Set.
func (g *Gulf) Task(name string, fn func() error, deps ...string) {
	g.s.Task(name, fn, deps...)
}

// Watch adds a set of patterns to be watched with corresponding tasks.
// There is no need to use Start.
func (g *Gulf) Watch(patterns []string, tasks ...string) {
	g.s.Watch(patterns, tasks...)
	g.watch = true
}

func main() {
	g := New()
	Tasks(g)
	err := g.s.Exec(os.Args[1])
	if g.watch {
		g.s.Start()
	}
	if err != nil {
		fmt.Println(err)
	}
}
