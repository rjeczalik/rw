package rw

import (
	"bytes"
	"io"
	"runtime"
	"strings"
	"testing"
)

func TestNewLineWriter(t *testing.T) {
	cases := [...]struct {
		s string
		w string
		u string
	}{
		{"dd\n", "dd\r\n", "dd\n"},
		{"", "", ""},
		{"\r\ni", "\r\ni", "\ni"},
		{"\ndd\n", "\r\ndd\r\n", "\ndd\n"},
		{"\n", "\r\n", "\n"},
		{"\r", "\r", "\r"},
		{"a\r\n", "a\r\n", "a\n"},
	}
	for i, cas := range cases {
		var buf bytes.Buffer
		wr := NewLineWriter(&buf)
		n, err := io.Copy(wr, strings.NewReader(cas.s))
		if err != nil {
			t.Errorf("want err=nil; got %v (i=%d)", err, i)
			continue
		}
		var w string
		if runtime.GOOS == "windows" {
			w = cas.w
		} else {
			w = cas.u
		}
		if got := buf.String(); got != w {
			t.Errorf("want got=%q; got %q (i=%d)", w, got, i)
			continue
		}
		if want := len(cas.s); int(n) != want {
			t.Errorf("want n=%d; got %d (i=%d)", want, n, i)
		}
	}
}
