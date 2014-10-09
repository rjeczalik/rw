package rw

import (
	"io/ioutil"
	"testing"
	"time"
)

func TestWaitWriter(t *testing.T) {
	cases := map[string][]string{
		"foo": {"opasd", "asdasdf", "osdfsdf", "rtyfoot"},
		"bar": {"sdfba", "sdfswe34", "w3psgb", "ar2al"},
		"baz": {"alqlqpqp102bBBAZ", "baz"},
		"qux": {"quxpwo300ls,"},
		"quu": {"dldl20slw", "wejznmqo10f8w0-ak", "04j0qm204n2", "ofg0qkquu"},
		"209skswkwprt,mxg;dfglsp": {"a", "b", "209skswk", "wprt,mx", "g;dfglsp"},
	}
	for msg, parts := range cases {
		w := WaitWriter(ioutil.Discard, []byte(msg))
		go func() {
			for _, s := range parts {
				if _, err := w.Write([]byte(s)); err != nil {
					panic(err)
				}
			}
		}()
		if err := w.Wait(time.Second); err != nil {
			t.Errorf("want w.Wait(...)=nil; got %v", err)
			continue
		}
	}
}
