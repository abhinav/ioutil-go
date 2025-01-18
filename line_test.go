package ioutil_test

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"

	"go.abhg.dev/io/ioutil"
)

func TestLineWriter(t *testing.T) {
	tests := []struct {
		name string
		give []string
		want []string
	}{
		{name: "Empty"},
		{
			name: "SingleLine",
			give: []string{"hello\n"},
			want: []string{"hello"},
		},
		{
			name: "SingleLineBuffered",
			give: []string{"hello"},
			want: []string{"hello"},
		},
		{
			name: "MultiLine",
			give: []string{"hello\n", "world\n"},
			want: []string{"hello", "world"},
		},
		{
			name: "MultiLineAcrossWrites",
			give: []string{"hello", "\n", "world\n"},
			want: []string{"hello", "world"},
		},
		{
			name: "LinesAcrossManyWrites",
			give: []string{"h", "el", "lo\nw", "or", "ld\n"},
			want: []string{"hello", "world"},
		},
		{
			name: "EmptyLine",
			give: []string{"foo\n", "\n", "bar\n"},
			want: []string{"foo", "", "bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			w, done := ioutil.LineWriter(func(bs []byte) {
				got = append(got, string(bs))
			})

			for _, s := range tt.give {
				if _, err := io.WriteString(w, s); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			done()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("want: %q\n got: %q", tt.want, got)
			}
		})
	}
}

func FuzzLineWriter(f *testing.F) {
	f.Add([]byte("hello\nworld\n"))
	f.Fuzz(func(t *testing.T, give []byte) {
		var gotBuff bytes.Buffer
		w, done := ioutil.LineWriter(func(bs []byte) {
			gotBuff.Write(bs)
			gotBuff.WriteByte('\n')
		})

		if _, err := w.Write(give); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		done()

		got := strings.TrimSuffix(gotBuff.String(), "\n")
		want := strings.TrimSuffix(string(give), "\n")
		if got != want {
			t.Errorf("want: %q\n got: %q", want, got)
		}
	})
}
