package ioutil_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"go.abhg.dev/io/ioutil"
)

func TestPrintfWriter(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		writes []string
		want   string
	}{
		{name: "Empty"},
		{
			name:   "SingleLine",
			writes: []string{"hello\n"},
			want:   "hello\n",
		},
		{
			name:   "Prefix",
			prefix: "prefix: ",
			writes: []string{"hello\nwo", "rld"},
			want: "prefix: hello\n" +
				"prefix: world\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.Buffer
			w, done := ioutil.PrintfWriter(func(format string, args ...any) {
				fmt.Fprintf(&got, format, args...)
				got.WriteByte('\n')
			}, tt.prefix)

			for _, s := range tt.writes {
				if _, err := io.WriteString(w, s); err != nil {
					t.Fatalf("Write failed: %v", err)
				}
			}
			done()

			if got := got.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
