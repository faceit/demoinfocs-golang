package bitread

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

type CaptureReader struct {
	*BitReader
	startRead bool
	stopRead  bool
	rdr       StoppableReader
	Out       *bytes.Buffer
}

func NewCaptureBitReader(underlying io.Reader) *CaptureReader {
	var buf bytes.Buffer
	rdr := NewStoppableReader(underlying, &buf)
	rdr.Begin()
	b := make([]byte, largeBuffer)
	br := newBitReader(rdr, &b)
	return &CaptureReader{BitReader: br, rdr: rdr, Out: &buf}
}

func (r *CaptureReader) BeginCapture() {
	// TODO this method will be called as we are midway through reading from the buffer
	// so we have to pickup everything that hasn't been read yet
	r.stopRead = false
	r.startRead = true
}

func (r *CaptureReader) BeginChunk(n int) {
	offset := r.ActualPosition() - r.LazyPosition()
	if r.startRead && !r.rdr.IsReading() && offset%8 == 0 {
		fmt.Printf("begin capture with buf size %d\n", r.Out.Len())
		r.Out.Write(r.Buffer[offset/8 : len(r.Buffer)-sled])
		fmt.Printf("filled from buffer, size is now %d\n", r.Out.Len())
		r.rdr.Begin()
	}
	r.BitReader.BeginChunk(n)
}

func (r *CaptureReader) WriteOut(filename string) error {
	return ioutil.WriteFile(filename, r.Out.Bytes(), 777)
}

func (r *CaptureReader) EndChunk() {
	r.BitReader.EndChunk()
	offset := r.ActualPosition() - r.LazyPosition()
	if r.stopRead && r.rdr.IsReading() && r.ChunkFinished() &&
		offset%8 == 0 {
		r.rdr.End()
		toClear := largeBuffer - sled - (offset / 8)
		r.Out.Truncate(r.Out.Len() - toClear)
		fmt.Printf("end capture with buf size %d\n", r.Out.Len())
	}
}

func (r *CaptureReader) EndCapture() {
	r.stopRead = true
	r.startRead = false
}
