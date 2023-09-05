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
	state            boundaryReaderState
}

//go:generate stringer -type=boundaryReaderState
type boundaryReaderState int

const (
	boundaryReaderStateEmpty = boundaryReaderState(iota)
	boundaryReaderStateNeedsBoundary
	boundaryReaderStateOnlyContent
	boundaryReaderStatePartialBoundaryInBuffer
	boundaryReaderStateCompleteBoundaryInBuffer
)

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
	if errors.IsEOF(err) {
	}

	if !errors.IsEOF(err) {
		return
	}

	br.resetRemainingContentIfNecessary()

	return
}

func (br *boundaryReader) resetRemainingContentIfNecessary() {
	switch br.state {
	default:
		panic(ErrInvalidBoundaryReaderState)

	case boundaryReaderStatePartialBoundaryInBuffer:
		eof := false
		br.remainingContent, eof = br.buffer.Find(br.boundary)

		if eof {
			br.state = boundaryReaderStatePartialBoundaryInBuffer
		} else {
			br.state = boundaryReaderStateCompleteBoundaryInBuffer
		}

	case boundaryReaderStateOnlyContent:
		// noop

	case boundaryReaderStateNeedsBoundary:
		// noop
	case boundaryReaderStateEmpty:
		// noop
	case boundaryReaderStateCompleteBoundaryInBuffer:
		// noop
	}

	return
}

func (br *boundaryReader) ReadBoundary() (n int, err error) {
	switch br.state {
	default:
		panic(ErrInvalidBoundaryReaderState)

	case boundaryReaderStateOnlyContent,
		boundaryReaderStateCompleteBoundaryInBuffer,
		boundaryReaderStatePartialBoundaryInBuffer:
		err = ErrExpectedContentRead
		return

	case boundaryReaderStateNeedsBoundary, boundaryReaderStateEmpty:
		// noop
	}

	eof := false
	n, eof = br.buffer.PeekMatchAdvance(br.boundary)

	switch {
	case n == len(br.boundary):
		br.state = boundaryReaderStatePartialBoundaryInBuffer

		var n1 int
		n1, err = br.fillBuffer()

		if errors.IsEOF(err) && n1 > 0 || br.buffer.Len() > 0 {
			// the buffer has more content, so try a boundary read again
			err = nil
		} else if err != nil {
			// the buffer has no more content, so this is an actual EOF
			return
		}

		return

	case n == br.buffer.Len() && eof:
		var n1 int
		n1, err = br.fillBuffer()

		if errors.IsEOF(err) && n1 > 0 && br.buffer.Len() > 0 {
			// the buffer has more content, so try a boundary read again
			err = nil
		} else if err != nil {
			// the buffer has no more content, so this is an actual EOF
			return
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
	switch br.state {
	default:
		panic(ErrInvalidBoundaryReaderState)

	case boundaryReaderStateNeedsBoundary,
		boundaryReaderStateEmpty:
		err = io.EOF
		return

	case boundaryReaderStateOnlyContent:
		n, err = br.buffer.Read(p)

		if errors.IsEOF(err) {
			err = nil

			_, err = br.fillBuffer()

			if err != nil && !errors.IsEOF(err) {
				return
			}
		}

	case boundaryReaderStatePartialBoundaryInBuffer:
		if len(p) > br.remainingContent {
			p = p[:br.remainingContent+1]
		}

		n, err = br.buffer.Read(p)

		if err != nil {
			panic("invalid state")
		}

		br.remainingContent -= n

		if br.remainingContent <= 0 {
			_, err = br.fillBuffer()

			if err != nil && !errors.IsEOF(err) {
				return
			}
		}

	case boundaryReaderStateCompleteBoundaryInBuffer:
		if len(p) > br.remainingContent {
			p = p[:br.remainingContent+1]
		}

		n, err = br.buffer.Read(p)

		if err != nil && !errors.IsEOF(err) {
			panic("invalid state")
		}

		br.remainingContent -= n

		if br.remainingContent <= 0 {
			br.state = boundaryReaderStateNeedsBoundary
		}
	}

	return
}
