package bitread

import (
	"fmt"
	"io"
	"time"
)

// TeeReader returns a Reader that writes to w what it reads from r.
// All reads from r performed through it are matched with
// corresponding writes to w. There is no internal buffering -
// the write must complete before the read completes.
// Any error encountered while writing is reported as a read error.
func NewStoppableReader(r io.Reader, w io.Writer) StoppableReader {
	return &teeReader{r: r, w: w}
}

type StoppableReader interface {
	io.Reader
	Begin()
	IsReading() bool
	End()
}

type teeReader struct {
	r    io.Reader
	w    io.Writer
	read bool
}

func (t teeReader) IsReading() bool {
	return t.read
}

func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 && t.read {
		if n, err := t.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}

func (t *teeReader) Begin() {
	if !t.read {
		fmt.Printf("started reading at %s\n", time.Now())
	}
	t.read = true
}

func (t *teeReader) End() {
	if t.read {
		fmt.Printf("stopped reading at %s\n", time.Now())
	}
	t.read = false
}
