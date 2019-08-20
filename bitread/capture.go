package bitread

import (
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
)

var pattern = regexp.MustCompile("Skipping data failed, expected position (\\d+) got (\\d+)")

type CaptureReader struct {
	*BitReader
	rdr StoppableReader
	Out bytes.Buffer
}

func NewCaptureBitReader(underlying io.Reader) *CaptureReader {
	var buf bytes.Buffer
	rdr := NewStoppableReader(underlying, &buf)
	rdr.Begin()

	br := bitReaderPool.Get().(*BitReader)
	br.Open(underlying, 32)
	return &CaptureReader{BitReader: br, rdr: rdr}
}

func (r *CaptureReader) BeginCapture() {
	r.rdr.Begin()
}

func (r *CaptureReader) WriteOut(filename string) error {
	return ioutil.WriteFile(filename, r.Out.Bytes(), 777)
}

func (r *CaptureReader) EndCapture() {
	defer func() {
		if rcv := recover(); rcv != nil {
			msg := rcv.(string)
			matches := pattern.FindStringSubmatch(msg)
			low, _ := strconv.Atoi(matches[1])
			high, _ := strconv.Atoi(matches[2])
			r.Skip(high - low)
			r.EndChunk()
		}
	}()
	if !r.ChunkFinished() {
		r.EndChunk()
	}

	r.rdr.End()
}
