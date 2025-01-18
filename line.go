package ioutil

import (
	"bytes"
	"io"
)

// LineWriter returns an io.Writer that will feed its writes
// through to the provided write function, one line at a time,
// not including the newline character.
//
// The done function must be called when no more writes are expected.
// This will flush any buffered writes to the write function.
//
// The returned io.Writer is NOT safe for concurrent use.
// Wrap it with [SyncWriter] to use it in a concurrent setting.
func LineWriter(write func([]byte)) (w io.Writer, done func()) {
	lw := &lineWriter{writeLine: write}
	return lw, lw.flush
}

// lineWriter is an io.Writer that writes to a log.Logger.
type lineWriter struct {
	writeLine func([]byte)
	buff      bytes.Buffer
}

var _newline = []byte{'\n'}

func (w *lineWriter) Write(bs []byte) (int, error) {
	total := len(bs)
	for len(bs) > 0 {
		var (
			line []byte
			ok   bool
		)
		line, bs, ok = bytes.Cut(bs, _newline)
		if !ok {
			// No newline. Buffer and wait for more.
			w.buff.Write(line)
			break
		}

		if w.buff.Len() == 0 {
			// No prior partial write. Flush.
			w.writeLine(line)
			continue
		}

		// Flush prior partial write.
		w.buff.Write(line)
		w.writeLine(w.buff.Bytes())
		w.buff.Reset()
	}
	return total, nil
}

// flush flushes buffered text, even if it doesn't end with a newline.
func (w *lineWriter) flush() {
	if w.buff.Len() > 0 {
		w.writeLine(w.buff.Bytes())
		w.buff.Reset()
	}
}
