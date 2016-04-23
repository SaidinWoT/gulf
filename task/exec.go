package task

import "sync"

// Exec runs the task with the provided name after resolving all of its dependencies.
// As long as the task exists, any error returned will be an ErrExec,
// which may be introspected for the errors returned by the dependencies.
func (s *Set) Exec(name string) error {
	if s.err != nil {
		return s.err
	}
	e := &exec{
		fs: make(map[string]func() error),
		ts: s.ts,
	}
	t, ok := s.ts[name]
	if !ok {
		return ErrTaskNotExist{name}
	}
	return e.runOnce(t)
}

type exec struct {
	sync.RWMutex
	fs map[string]func() error
	ts map[string]task
}

func (e *exec) run(t task) error {
	errs := newErrExec()
	wg := new(sync.WaitGroup)
	wg.Add(len(t.deps))
	for _, dep := range t.deps {
		go func(d string) {
			var err error
			f, d := parseFlags(d)
			dt, ok := e.ts[d]
			if !ok {
				err = ErrTaskNotExist{name: d}
			} else if dt.f.multi || f.multi {
				err = e.run(dt)
			} else {
				err = e.runOnce(dt)
			}
			if err != nil {
				errs.Add(d, err, dt.f.optional || f.optional)
			}
			wg.Done()
		}(dep)
	}
	wg.Wait()
	if errs.Failed() {
		return errs
	}
	errs.Task = t.fn()
	if errs.Task != nil {
		return errs
	}
	return nil
}

func (e *exec) runOnce(t task) error {
	e.RLock()
	fn, ok := e.fs[t.name]
	e.RUnlock()
	if ok {
		return fn()
	}
	e.Lock()
	fn, ok = e.fs[t.name]
	if !ok {
		var (
			err error
			o   sync.Once
		)
		fn = func() error {
			o.Do(func() {
				err = e.run(t)
			})
			return err
		}
		e.fs[t.name] = fn
	}
	e.Unlock()
	return fn()
}
