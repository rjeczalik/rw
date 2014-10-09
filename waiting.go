package rw

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

// WaitingWriter TODO
type WaitingWriter struct {
	W io.Writer // underlying Writer
	P []byte    // byte sequence to wait for

	spin uint32
	offs int
}

// WaitWriter TODO
func WaitWriter(w io.Writer, p []byte) *WaitingWriter {
	return &WaitingWriter{
		W: w,
		P: p,
	}
}

// Write TODO
func (ww *WaitingWriter) Write(p []byte) (int, error) {
	for _, b := range p {
		if ww.P[ww.offs] == b {
			ww.offs++
			if ww.offs == len(ww.P) {
				ww.offs = 0
				atomic.CompareAndSwapUint32(&ww.spin, 0, 1)
			}
		} else {
			ww.offs = 0
		}
	}
	return ww.W.Write(p)
}

// Wait TODO
func (ww *WaitingWriter) Wait(d time.Duration) error {
	t := time.After(d)
	for {
		select {
		case <-t:
			return fmt.Errorf("timeout waiting for p=%q", ww.P)
		default:
			if atomic.CompareAndSwapUint32(&ww.spin, 1, 0) {
				return nil
			}
		}
	}
}
