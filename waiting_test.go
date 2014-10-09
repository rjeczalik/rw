package rw

import (
	"io"
	"io/ioutil"
	"testing"
	"time"
)

var cases = map[string][]string{
	"foo": {"opasd", "asdasdf", "osdfsdf", "rtyfoot"},
	"bar": {"sdfba", "sdfswe34", "w3psgb", "ar2al"},
	"baz": {"alqlqpqp102bBBAZ", "baz"},
	"qux": {"quxpwo300ls,"},
	"quu": {"dldl20slw", "wejznmqo10f8w0-ak", "04j0qm204n2", "ofg0qkquu"},
	"209skswkwprt,mxg;dfglsp": {"a", "b", "209skswk", "wprt,mx", "g;dfglsp"},
}

func write(w io.Writer, s []string) {
	for _, s := range s {
		if _, err := w.Write([]byte(s)); err != nil {
			panic(err)
		}
	}
}

func TestWaitWriter(t *testing.T) {
	for msg, parts := range cases {
		w := WaitWriter(ioutil.Discard, []byte(msg))
		go write(w, parts)
		if err := w.Wait(time.Second); err != nil {
			t.Errorf("want w.Wait(...)=nil; got %v", err)
			continue
		}
	}
}

func TestWaitWriterTimeout(t *testing.T) {
	for msg, parts := range cases {
		w := WaitWriter(ioutil.Discard, []byte(msg+"must timeout"))
		err := make(chan error)
		go write(w, parts)
		go func() { err <- w.Wait(20 * time.Millisecond) }()
		select {
		case err := <-err:
			if err == nil {
				t.Errorf("want err != nil (msg=%s)", msg)
			}
		case <-time.After(time.Second):
			t.Errorf("test has timed out (msg=%s)", msg)
		}
	}
}
