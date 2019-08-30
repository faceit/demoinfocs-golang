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
	r := &CaptureReader{BitReader: br, rdr: rdr, Out: &buf}
	r.startRead = true
	return r
}

func (r *CaptureReader) BeginCapture() {
	r.stopRead = false
	r.startRead = true
}

func (r *CaptureReader) WriteOut(filename string) error {
	return ioutil.WriteFile(filename, r.Out.Bytes(), 777)
}

func (r *CaptureReader) EndChunk() {
	r.BitReader.EndChunk()
	offset := r.ActualPosition() - r.LazyPosition()
	atBoundary := offset%8 == 0 && len(r.ChunkTargets) == 0

	if r.rdr.IsReading() && r.stopRead && atBoundary {
		r.rdr.End()
		toClear := largeBuffer - (offset / 8)
		r.Out.Truncate(r.Out.Len() - toClear)
		fmt.Printf("end capture with buf size %d\n", r.Out.Len())
	}

	if !r.rdr.IsReading() && r.startRead && atBoundary {
		fmt.Printf("begin capture with buf size %d\n", r.Out.Len())
		r.Out.Write(r.Buffer[offset/8:])
		fmt.Printf("filled from buffer, size is now %d\n", r.Out.Len())
		r.rdr.Begin()
	}
}

func (r *CaptureReader) EndCapture() {
	r.stopRead = true
	r.startRead = false
}
