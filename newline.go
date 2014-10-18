package rw

import (
	"bytes"
	"io"
	"runtime"
	"strings"
)

type newLineWriter struct {
	W io.Writer
	f newLineFunc
}

// NewLineWriter returns writer wrapping provided io.Writer and updating
// passed data by using proper new lines for specific platform.
func NewLineWriter(w io.Writer) io.Writer {
	wr := &newLineWriter{
		W: w,
	}
	if runtime.GOOS == "windows" {
		wr.f = winNewLine
	} else {
		wr.f = unixNewLine
	}
	return wr
}

type newLineFunc func(io.Writer, []byte) (int, error)

// winNewLine implements writing for windows.
func winNewLine(w io.Writer, p []byte) (n int, err error) {
	s, pos := string(p), 0
	var b bytes.Buffer
	var ss string
	for i := strings.Index(s, "\n"); i != -1; i = strings.Index(s[pos:], "\n") {
		if _, err = b.WriteString(s[pos : pos+i]); err != nil {
			return 0, err
		}
		if (i > 0 && s[pos+i-1] != '\r') || i == 0 {
			ss = "\r\n"
		} else {
			ss = "\n"
		}
		if _, err = b.WriteString(ss); err != nil {
			return 0, err
		}
		pos += i + 1
	}
	if _, err = b.WriteString(s[pos:]); err != nil {
		return 0, err
	}
	if n, err = w.Write(b.Bytes()); err != nil {
		return
	}
	return len(p), nil
}

// unixNewLine implements writing for unix.
func unixNewLine(w io.Writer, p []byte) (n int, err error) {
	s := strings.Replace(string(p), "\r\n", "\n", -1)
	if n, err = w.Write([]byte(s)); err == nil {
		return len(p), nil
	}
	return
}

// Write implements io.Writer.
func (nl *newLineWriter) Write(p []byte) (n int, err error) {
	return nl.f(nl.W, p)
}
