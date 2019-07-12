package main

import "io"

// TeeReader returns a Reader that writes to w what it reads from r.
// All reads from r performed through it are matched with
// corresponding writes to w. There is no internal buffering -
// the write must complete before the read completes.
// Any error encountered while writing is reported as a read error.
func MagicReader(r io.Reader, w io.Writer) io.Reader {
	return &teeReader{r: r, w: w}
}

type teeReader struct {
	r    io.Reader
	w    io.Writer
	read bool
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
	t.read = true
}

func (t *teeReader) End() {
	t.read = false
}
