package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/rjeczalik/rw"
)

const usage = `usage: fsplit [-limit SIZE_IN_MB] FILE`

var limit = flag.Int("limit", 50, "Single file part size limit in megabytes.")

func die(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		die(usage)
	}

	if err := fsplit(flag.Arg(0), *limit*rw.MB); err != nil {
		die(err)
	}
}

func fsplit(path string, limit int) error {
	r := rw.LimitReader(path)
	w := rw.LimitWriter(path, limit)
	_, err := io.Copy(w, r)
	return nonil(err, r.Close(), w.Close())
}

func nonil(err ...error) error {
	for _, e := range err {
		if e != nil {
			return e
		}
	}
	return nil
}
