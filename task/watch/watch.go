package watch

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/SaidinWoT/gulf/glob"
	"gopkg.in/fsnotify.v1"
)

type timers struct {
	sync.RWMutex
	m     map[string]*time.Timer
	delay time.Duration
	exec  func(string) error
}

// Watch adds a set of tasks to the list that will be executed when any of the provided patterns are matched.
func (s *Set) Watch(patterns []string, tasks ...string) {
	for _, p := range patterns {
		s.wm[p] = append(s.wm[p], tasks...)
	}
}

// Start begins watching all of the directories containing files added to s.
func (s *Set) Start() error {
	w, err := s.dirWatcher()
	if err != nil {
		return err
	}

	ts := &timers{
		m:     make(map[string]*time.Timer),
		delay: s.delay,
		exec:  s.Exec,
	}

	for {
		select {
		case event := <-w.Events:
			for p, tasks := range s.wm {
				if matches, _ := s.match(p, event.Name); matches {
					for _, task := range tasks {
						go ts.debounce(task)
					}
				}
			}
		case _ = <-w.Errors:
		}
	}
}

func (s *Set) dirWatcher() (*fsnotify.Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	var paths []string
	for path := range s.wm {
		paths = append(paths, path)
	}
	dirs := make(map[string]struct{})
	for _, path := range glob.Parse(s.glob, paths...) {
		dirs[filepath.Dir(path)] = struct{}{}
	}
	for dir := range dirs {
		w.Add(dir)
	}
	return w, nil
}

// debounce creates an AfterFunc to execute the provided task
// or postpones its execution if the AfterFunc already exists.
// This behavior is intended to deal with the series of events
// that many popular text editors issue when writing a file.
func (ts *timers) debounce(name string) {
	ts.RLock()
	t, ok := ts.m[name]
	ts.RUnlock()
	if ok {
		t.Reset(ts.delay)
		return
	}
	ts.Lock()
	ts.m[name] = time.AfterFunc(ts.delay, func() {
		ts.Lock()
		delete(ts.m, name)
		ts.Unlock()
		ts.exec(name)
	})
	ts.Unlock()
}
