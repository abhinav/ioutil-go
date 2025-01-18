package ioutil

import (
	"io"
	"sync"
)

// SyncWriter makes an io.Writer safe for concurrent use
// by adding a mutex around the Write calls.
func SyncWriter(w io.Writer) io.Writer {
	return &syncWriter{w: w}
}

type syncWriter struct {
	w io.Writer
	m sync.Mutex
}

func (w *syncWriter) Write(p []byte) (n int, err error) {
	w.m.Lock()
	defer w.m.Unlock()
	return w.w.Write(p)
}
