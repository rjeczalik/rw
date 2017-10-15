package rw

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
	PB = 1024 * TB
)

// BUG(rjeczalik): LimitFile creates a writer that always truncates
// the file if it exists - the append support is missing.

// NOTE: limitedFile's writer is always doing copy-on-write - it
// writes to a temporary files that are renamed on a successful
// Close() call.
//
// NOTE: It is not safe to call Read() and/or Write() methods
// concurrently - synchronization is left up to the caller.
type limitedFile struct {
	path    string
	n       int
	written int
	limit   int
	r       *os.File
	w       *os.File
	err     error
	tmp     string
	cow     []string
	single  bool
}

var _ io.ReadWriteCloser = (*limitedFile)(nil)

func LimitReader(path string) io.ReadCloser {
	return &limitedFile{
		path: path,
	}
}

func LimitWriter(path string, limit int) io.WriteCloser {
	return &limitedFile{
		path:  path,
		limit: limit,
	}
}

func (f *limitedFile) Read(p []byte) (int, error) {
	if f.r == nil && f.n == 0 {
		var err error
		switch f.r, err = os.Open(f.path); {
		case os.IsNotExist(err):
			var e error
			if f.r, e = os.Open(f.path + ".1"); e != nil {
				return 0, err // return original error
			}
			f.n++
		case err == nil:
			f.single = true
		default:
			return 0, err
		}
	}
	n, err := f.r.Read(p)
	if err == io.EOF {
		_ = f.r.Close()
		var e error
		if f.r, e = os.Open(fmt.Sprintf("%s.%d", f.path, f.n+1)); e != nil {
			if os.IsNotExist(e) {
				return n, err
			}
			return n, e
		}
		f.n++
		err = nil
	}
	return n, err
}

func (f *limitedFile) mktmp() (*os.File, error) {
	return ioutil.TempFile(filepath.Split(f.path))
}

func (f *limitedFile) Write(p []byte) (int, error) {
	if f.w == nil && f.n == 0 {
		var err error
		if f.w, err = f.mktmp(); err != nil {
			return 0, err
		}
	}
	n, err := f.w.Write(p)
	if err != nil {
		return n, err
	}
	f.written += n
	// Best-effort attempt of rotating before we approach limit.
	//
	// TODO(rjeczalik): If we wanted to be accurate, we could split
	// the write into two, if the write would exceed the limit.
	if f.written >= (f.limit - n) {
		if err = f.w.Close(); err != nil {
			return n, err
		}
		f.cow = append(f.cow, f.w.Name())
		if f.w, err = f.mktmp(); err != nil {
			return n, err
		}
		f.written = 0
	}
	return n, nil
}

func (f *limitedFile) Close() error {
	var err error
	if f.w != nil {
		err = nonil(err, f.w.Close())
		f.cow = append(f.cow, f.w.Name())
	}
	if f.r != nil {
		err = nonil(err, f.r.Close())
	}
	if err != nil {
		return err
	}
	if len(f.cow) == 1 {
		err = nonil(err, os.Rename(f.cow[0], f.path))
	} else {
		for i, tmp := range f.cow {
			path := fmt.Sprintf("%s.%d", f.path, i+1)
			err = nonil(err, os.Rename(tmp, path))
		}
		if f.single && err == nil {
			err = nonil(err, os.Remove(f.path))
		}
	}
	return err
}
