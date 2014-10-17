package rw

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestNewLineWriterWindows(t *testing.T) {
	cases := [...]struct {
		s string
		r string
	}{
		{"dd\n", "dd\r\n"},
		{"", ""},
		{"\r\ni", "\r\ni"},
		{"\ndd\n", "\r\ndd\r\n"},
		{"\n", "\r\n"},
		{"\r", "\r"},
		{"a\r\n", "a\r\n"},
	}
	for i, cas := range cases {
		var buf bytes.Buffer
		wr := &newLineWriter{&buf, winNewLine}
		n, err := io.Copy(wr, strings.NewReader(cas.s))
		if err != nil {
			t.Errorf("want err=nil; got %v (i=%d)", err, i)
			continue
		}
		if got := buf.String(); got != cas.r {
			t.Errorf("want got=%q; got %q (i=%d)", cas.r, got, i)
			continue
		}
		if want := len(cas.s); int(n) != want {
			t.Errorf("want n=%d; got %d (i=%d)", want, n, i)
		}
	}
}

func TestNewLineWriterUnix(t *testing.T) {
	cases := [...]struct {
		s string
		r string
	}{
		{"dd\n", "dd\n"},
		{"", ""},
		{"\r\ni", "\ni"},
		{"\ndd\n", "\ndd\n"},
		{"\n", "\n"},
		{"\r", "\r"},
		{"a\r\n", "a\n"},
	}
	for i, cas := range cases {
		var buf bytes.Buffer
		wr := &newLineWriter{&buf, unixNewLine}
		n, err := io.Copy(wr, strings.NewReader(cas.s))
		if err != nil {
			t.Errorf("want err=nil; got %v (i=%d)", err, i)
			continue
		}
		if got := buf.String(); got != cas.r {
			t.Errorf("want got=%q; got %q (i=%d)", cas.r, got, i)
			continue
		}
		if want := len(cas.s); int(n) != want {
			t.Errorf("want n=%d; got %d (i=%d)", want, n, i)
		}
	}
}
