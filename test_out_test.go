package ioutil_test

import (
	"fmt"
	"reflect"
	"testing"

	"go.abhg.dev/io/ioutil"
)

func TestTestLogWriter(t *testing.T) {
	var fakeT fakeTestLogger
	w := ioutil.TestLogWriter(&fakeT, "prefix: ")
	fmt.Fprint(w, "hello\nworld")

	// "world" is buffered, so it won't show yet.
	if got, want := fakeT.Logs(), []string{"prefix: hello"}; !reflect.DeepEqual(got, want) {
		t.Errorf("want: %q\n got: %q", want, got)
	}

	fakeT.RunCleanup()

	if got, want := fakeT.Logs(), []string{
		"prefix: hello",
		"prefix: world",
	}; !reflect.DeepEqual(got, want) {
		t.Errorf("want: %q\n got: %q", want, got)
	}
}

type fakeTestLogger struct {
	logs    []string
	cleanup func()
}

func (f *fakeTestLogger) Logs() []string { return f.logs }

func (f *fakeTestLogger) RunCleanup() {
	if f.cleanup != nil {
		f.cleanup()
	}
}

func (f *fakeTestLogger) Logf(format string, args ...any) {
	f.logs = append(f.logs, fmt.Sprintf(format, args...))
}

func (f *fakeTestLogger) Cleanup(fn func()) {
	old := f.cleanup
	f.cleanup = func() {
		fn()
		if old != nil {
			old()
		}
	}
}
