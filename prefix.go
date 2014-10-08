package rw

import (
	"bytes"
	"io"
)

// PrefixedWriter TODO
type PrefixedWriter struct {
	W      io.Writer     // underlying writer
	Prefix func() string // generates string used to prefix each line
}

// PrefixWriter TODO
func PrefixWriter(writer io.Writer, prefix func() string) io.Writer {
	return PrefixedWriter{W: writer, Prefix: prefix}
}

// Write TODO
func (pw PrefixedWriter) Write(p []byte) (int, error) {
	var (
		i   int
		err error
		buf bytes.Buffer
	)
	// If p contains multiple newlines we loop over each of them.
	for j, n := indexnl(p); j != -1; j, n = indexnl(p[i:]) {
		j += i
		// Write prefix.
		if _, err = buf.WriteString(pw.Prefix()); err != nil {
			return 0, err
		}
		// Write line.
		if _, err = buf.Write(p[i : j+n]); err != nil {
			return 0, err
		}
		i = j + n
	}
	if i != 0 {
		// Write last line if p does not end with a newline.
		if i < len(p) {
			// Write prefix for the last line.
			if _, err = buf.WriteString(pw.Prefix()); err != nil {
				return 0, err
			}
			// Write the last line.
			if _, err = buf.Write(p[i:]); err != nil {
				return 0, err
			}
		}
		n, err := io.Copy(pw.W, &buf)
		return min(int(n), len(p)), err
	}
	return pw.W.Write(p)
}
