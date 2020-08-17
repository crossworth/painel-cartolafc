package multiwriter

import (
	"errors"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

// MultiWriter is a synced multiwriter
// it will call sync if the writer is a file
type MultiWriter struct {
	lock sync.Mutex
	IO1  io.Writer
	IO2  io.Writer
}

// Write write to both registered writers and if its a file it will
// call sync on it to make sure it's flushed
func (mw *MultiWriter) Write(p []byte) (int, error) {
	if mw.IO1 == nil && mw.IO2 == nil {
		return 0, errors.New("no writer provided")
	}

	mw.lock.Lock()
	defer mw.lock.Unlock()

	var n1, n2 int
	var err1, err2 error

	if mw.IO1 != nil {
		n1, err1 = mw.IO1.Write(p)

		if file, ok := mw.IO1.(*os.File); ok {
			_ = file.Sync()
		}
	}

	if mw.IO2 != nil {
		n2, err2 = mw.IO2.Write(p)

		if file, ok := mw.IO2.(*os.File); ok {
			_ = file.Sync()
		}
	}

	n := n2
	if n1 > n2 {
		n = n1
	}

	err := err2
	if err1 != nil {
		err = err1
	}

	// NOTE(Pedro): On Windows the terminal output will be truncated sometimes
	// we can avoid it by sleeping for one Nanosecond
	if runtime.GOOS == "windows" {
		time.Sleep(1 * time.Nanosecond)
	}

	return n, err
}
