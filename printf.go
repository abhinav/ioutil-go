package ioutil

import (
	"io"
)

// PrintfWriter builds an io.Writer that will feed its writes
// through to the provided printf-style function, one line at a time,
// not including the newline character.
//
// The provided prefix will be prepended to each line written.
// Leave prefix empty to write lines as is.
//
// The done function must be called when no more writes are expected.
// This will flush any buffered writes to the write function.
//
// Use this to redirect the output of an io.Writer into a logger,
// or other printf-style function. For example:
//
//	w, done := ioutil.PrintfWriter(myLogger.Printf)
//	defer done()
//
//	cmd := exec.Command("some", "--long", "--running", "command")
//	cmd.Stdout = w
//	cmd.Stderr = w
//	err := cmd.Run()
//
// The returned io.Writer is NOT safe for concurrent use.
// Wrap it with [SyncWriter] to use it in a concurrent setting.
func PrintfWriter(printf func(string, ...any), prefix string) (w io.Writer, done func()) {
	return LineWriter(func(bs []byte) {
		printf("%s%s", prefix, bs)
	})
}
