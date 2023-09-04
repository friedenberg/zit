package ohio

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type BoundaryReader interface {
	io.Reader
	ReadBoundary() (int, error)
}

type boundaryReader struct {
	reader           *bufio.Reader
	boundary         []byte
	remainingContent int
	buffer           *RingBuffer
}

func MakeBoundaryReader(r io.Reader, boundary string) BoundaryReader {
	// TODO-P1 perf allow for optimized buffer size
	d := 0

	if len(boundary) > ringBufferDefaultSize {
		d = len(boundary)
	}

	b := MakeRingBuffer(d)

	return &boundaryReader{
		reader:   bufio.NewReader(r),
		boundary: []byte(boundary),
		buffer:   b,
	}
}

func (br *boundaryReader) fillBuffer() (n int, err error) {
	n, err = br.buffer.ReadFromSmall(br.reader)

	return
}

func (br *boundaryReader) ReadBoundary() (n int, err error) {
	if br.remainingContent > 0 {
		err = ErrExpectedContentRead
		return
	}

	eof := false
	n, eof = br.buffer.PeekMatchAdvance(br.boundary)

	switch {
	case n == len(br.boundary):
		br.remainingContent, _ = br.buffer.Find(br.boundary)
		return

	case n == br.buffer.Len() && eof:
		var n1 int
		n1, err = br.fillBuffer()

		if err != nil {
			if errors.IsEOF(err) && n1 > 0 {
				// the buffer has more content, so try a boundary read again
				err = nil
			} else {
				// the buffer has no more content, so this is an actual EOF
				return
			}
		}
		// read more and try again
		return br.ReadBoundary()

	case n < len(br.boundary):
		fallthrough

	default:
		// not a match, fail
		err = ErrBoundaryNotFound
		return
	}
}

func (br *boundaryReader) Read(p []byte) (n int, err error) {
	if br.remainingContent == -1 {
		err = io.EOF
		return
	}

	if len(p) < br.buffer.Len() {
		var n1 int
		n1, err = br.fillBuffer()

		if err != nil {
			if errors.IsEOF(err) && n1 > 0 {
				// the buffer has more content
				err = nil
			} else {
				// the buffer has no more content, so this is an actual EOF
				return
			}
		}
	}

	if len(p) > br.remainingContent {
		p = p[:br.remainingContent+1]
	}

	n, err = br.buffer.Read(p)
	br.remainingContent -= n

	if err == nil && br.remainingContent == 0 {
		err = io.EOF
	}

	return
}
