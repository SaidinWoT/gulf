# gulf
Concurrent streaming build system in Go.

# Installation
To install the command, `go get -u github.com/SaidinWoT/gulf/cmd/gulf` followed
by `go install github.com/SaidinWoT/gulf/cmd/gulf`.  If you're interested in
using the innards, see the directories for READMEs on the various packages that
make up this project.

# Usage
Users familiar with [gulpjs](https://github.com/gulpjs/gulp) will find this
largely familiar, if more imperative.

Once the command is installed, simply create a gulf.go file with the following
skeleton:
```go
// +build gulf

package main

func Tasks(g *Gulf) {
}
```
and define a set of tasks of the form `g.Task("name", func() error { /* do some
things */ }, "optional", "dependencies", "here")`.  With that file in place,
run `gulf name` whenever you want to run the task and let gulf take care of the
rest.

# Documentation
The Gulf type provides a wrapper around `gulf/stream` and `gulf/task/watch`
for convenience.  Documentation is available on the
[gulf/cmd](https://godoc.org/github.com/SaidinWoT/gulf/cmd/) package, which
exists only to document what is available for use in gulf.go.
