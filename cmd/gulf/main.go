package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

//go:generate strmirror boilerplate.go
const boilerplateLoc = "/src/github.com/SaidinWoT/gulf/cmd/gulf/boilerplate.go"

var (
	runFlag   bool
	buildFlag bool
)

func init() {
	flag.BoolVar(&runFlag, "r", false, "Run a task from gulf.go without building the binary.")
	flag.BoolVar(&buildFlag, "b", false, "Rebuild the binary regardless of gulf.go's modtime.")
}

func main() {
	flag.Parse()
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	if runFlag {
		run(wd, os.Args...)
		return
	}
	binTime := modTime(filepath.Join(wd, "gulf"))
	rawTime := modTime(filepath.Join(wd, "gulf.go"))
	if rawTime.IsZero() {
		fmt.Println("There was an error accessing gulf.go. This is problematic.")
		return
	}
	if rawTime.After(binTime) || buildFlag {
		fmt.Println("Rebuilding the local gulf binary.")
		rebuild(wd)
	}
	task := "default"
	if len(flag.Args()) > 0 {
		task = flag.Arg(0)
	}
	fmt.Println("Running task", task)
	runCmd(filepath.Join(wd, "gulf"), task)
}

func run(wd string, args ...string) error {
	var s string
	var err error
	if b := boilerplate(); b == "" {
		f, err := ioutil.TempFile(wd, "")
		if err != nil {
			return err
		}
		f.WriteString(boilerplate_string)
		s = f.Name()
	} else {
		s, err = safeName(wd)
		if err != nil {
			return err
		}
		os.Link(b, s)
	}
	defer os.Remove(s)
	return runCmd("go", "run", "-tags=gulf", "gulf.go", s, flag.Arg(0))
}

func rebuild(wd string) error {
	dir, err := ioutil.TempDir(wd, "")
	if err != nil {
		fmt.Println(err)
		return err
	}
	if b := boilerplate(); b == "" {
		writeBoilerplate(dir)
	} else {
		os.Link(b, filepath.Join(dir, "boilerplate.go"))
	}
	os.Link(filepath.Join(wd, "gulf.go"), filepath.Join(dir, "gulf.go"))
	os.Chdir(dir)
	defer func() {
		os.Chdir(wd)
		os.RemoveAll(dir)
	}()
	return runCmd("go", "build", "-tags=gulf", "-o", filepath.Join(wd, "gulf"))
}

func boilerplate() string {
	loc := os.ExpandEnv("$GULFCMDSRC")
	if loc != "" {
		if _, err := os.Stat(loc); err == nil {
			return filepath.Join(loc, "boilerplate.go")
		}
	}
	gopath := os.ExpandEnv("$GOPATH")
	paths := filepath.SplitList(gopath)
	if len(paths) == 0 {
		return ""
	}
	for _, p := range paths {
		loc := filepath.Join(p, boilerplateLoc)
		_, err := os.Stat(loc)
		if err == nil {
			return loc
		}
	}
	return ""
}

func writeBoilerplate(dir string) error {
	f, err := os.Create(filepath.Join(dir, "boilerplate.go"))
	if err == nil {
		_, err = f.WriteString(boilerplate_string)
	}
	return err
}

func modTime(path string) time.Time {
	fi, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return fi.ModTime()
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}
