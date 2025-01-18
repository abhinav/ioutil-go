package ioutil

import (
	"io"
)

// TestLogger is a subset of the testing.TB interface
// with support for logging messages and cleaning up resources.
type TestLogger interface {
	Logf(string, ...any)
	Cleanup(func())
}

// TestLogWriter builds an io.Writer that will feed its writes
// through to the given test log.
//
// Prefix is prepended to each line written.
// Leave prefix empty to write lines as is.
//
// The returned io.Writer is NOT safe for concurrent use.
// Wrap it with [SyncWriter] to use it in a concurrent setting.
func TestLogWriter(t TestLogger, prefix string) io.Writer {
	w, done := PrintfWriter(t.Logf, prefix)
	t.Cleanup(done)
	return w
}
