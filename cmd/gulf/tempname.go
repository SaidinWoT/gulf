package main

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// This is all stolen from io/ioutil's tmepfile.go.
// We simply want a name that doesn't exist without creating a file
var (
	rand   uint32
	randmu sync.Mutex
)

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func next() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

func safeName(dir string) (s string, err error) {
	nconflict := 0
	for i := 0; i < 10000; i++ {
		s = filepath.Join(dir, next()+".go")
		_, err = os.Stat(s)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				rand = reseed()
			}
			continue
		}
		err = nil
		break
	}
	return
}
