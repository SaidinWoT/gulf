// This file exists simply to provide a nice godoc for the Gulf struct available in gulf.go.
// Importing this package would be a thoroughly useless endeavor.
package cmd

import (
	"time"

	"github.com/SaidinWoT/gulf/stream"
)

// Gulf is a simple struct to bring together gulf's functionality.
type Gulf struct {
	unexported struct{}
}

// The Option type provides a functional means of modifying a Gulf struct.
type Option func(*Gulf) error

// SetOption modifies g with the provided Options.
func (g *Gulf) SetOption(opts ...Option) error {
	return nil
}

func nopOption(_ *Gulf) error {
	return nil
}

// Globber returns an Option that sets a Gulf's globbing function to fn.
//
// Default: gulf/glob.Glob
func Globber(fn func(string) ([]string, error)) Option {
	return nopOption
}

// Matcher returns an Option that sets a Gulf's matcher function fn.
//
// Default: gulf/glob.Match
func Matcher(fn func(string, string) (bool, error)) Option {
	return nopOption
}

// WatchDelay returns an Option that sets the delay used to detect unique file events.
//
// Default: 10 * time.Millisecond
func Delay(d time.Duration) Option {
	return nopOption
}

// Src provides a simple wrapper around stream.Src.
// It parses patterns with g's globbing function and provides those to stream.Src.
//
// By default, a Gulf uses gulf/glob's Glob function.
// The behavior differs from filepath.Glob as follows:
// Globstar (**) matches 0 or more directories.
// Globs (*) only include dotfiles if there is an explicit dot before the glob character.
func (g *Gulf) Src(patterns ...string) stream.Stream {
	return stream.Stream{}
}

// Task adds a task to g's task Set.
func (g *Gulf) Task(name string, fn func() error, deps ...string) {}

// Watch adds a set of patterns to be watched with corresponding tasks.
// There is no need to use Start.
func (g *Gulf) Watch(patterns []string, tasks ...string) {}
