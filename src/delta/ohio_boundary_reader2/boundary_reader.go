package ohio_boundary_reader

import (
	"bufio"
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/ohio_ring_buffer2"
)

type BoundaryReader interface {
	io.Reader
	ReadBoundary() (int, error)
}

type boundaryReader struct {
	reader           *bufio.Reader
	ff               ohio_ring_buffer2.FindFunc
	remainingContent int
	lastFindIndex    int64
	buffer           *ohio_ring_buffer2.RingBuffer
	state            boundaryReaderState
}

func (br *boundaryReader) Reset(r io.Reader) {
	br.reader.Reset(r)
	br.remainingContent = 0
	br.buffer.Reset()
	br.state = boundaryReaderStateEmpty
}

func MakeBoundaryReader(r io.Reader, boundary string) BoundaryReader {
	d := 0

	if len(boundary) > ohio_ring_buffer2.RingBufferDefaultSize {
		d = len(boundary)
	}

	b := ohio_ring_buffer2.MakeRingBuffer(d)

	return &boundaryReader{
		reader: bufio.NewReader(r),
		buffer: b,
		ff:     ohio_ring_buffer2.FindBoundary([]byte(boundary)),
	}
}

func MakeBoundaryReaderPageSize(
	r io.Reader,
	boundary string,
	size int,
) BoundaryReader {
	if len(boundary) > ohio_ring_buffer2.RingBufferDefaultSize {
		size = len(boundary)
	}

	b := ohio_ring_buffer2.MakeRingBuffer(size)

	return &boundaryReader{
		reader: bufio.NewReader(r),
		buffer: b,
		ff:     ohio_ring_buffer2.FindBoundary([]byte(boundary)),
	}
}

func (br *boundaryReader) setState(s boundaryReaderState) {
	br.state = s
}

func (br *boundaryReader) fillBuffer() (n int, err error) {
	n, err = br.buffer.FillWith(br.reader)

	if err != nil && err != io.EOF {
		return
	}

	br.resetRemainingContentIfNecessary()

	return
}

func (br *boundaryReader) resetRemainingContentIfNecessary() {
	switch br.state {
	default:
		panic(ErrInvalidBoundaryReaderState)

	case boundaryReaderStatePartialBoundaryInBuffer,
		boundaryReaderStateOnlyContent:
		untilBoundary, _, partial := br.buffer.FindAnywhere(br.ff)

		switch {
		case untilBoundary > br.buffer.Len():
			panic("invalid state")

		case untilBoundary == -1 || partial:
			br.remainingContent = br.buffer.Len()

			if br.remainingContent == 0 {
				br.setState(boundaryReaderStateNeedsBoundary)
			} else {
				br.setState(boundaryReaderStateOnlyContent)
			}

		case partial:
			br.setState(boundaryReaderStatePartialBoundaryInBuffer)
			fallthrough

		case !partial:
			br.setState(boundaryReaderStateCompleteBoundaryInBuffer)
			fallthrough

		default:
			br.remainingContent = untilBoundary
		}

	case boundaryReaderStateNeedsBoundary:
		// noop
	case boundaryReaderStateEmpty:
		// noop
	case boundaryReaderStateCompleteBoundaryInBuffer:
		// noop
	}
}

func (br *boundaryReader) ReadBoundary() (length int, err error) {
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

	partial := false
	length, partial = br.buffer.FindFromStartAndAdvance(br.ff)

	switch {
	case length != 0 && !partial:
		br.setState(boundaryReaderStatePartialBoundaryInBuffer)

		var bytesReadFromRingBuffer int
		bytesReadFromRingBuffer, err = br.fillBuffer()

		if err == io.EOF && bytesReadFromRingBuffer > 0 || br.buffer.Len() > 0 {
			// the buffer has more content, so try a boundary read again
			err = nil
		} else if err != nil {
			// the buffer has no more content, so this is an actual EOF
			return
		}

		return

	case (length == br.buffer.Len() && partial) || br.buffer.Len() == 0:
		var n1 int
		n1, err = br.fillBuffer()

		if err == io.EOF && n1 > 0 && br.buffer.Len() > 0 {
			// the buffer has more content, so try a boundary read again
			err = nil
		} else if err != nil {
			// the buffer has no more content, so this is an actual EOF
			return
		}

		// read more and try again
		return br.ReadBoundary()

	case partial:
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
		br.remainingContent -= n

		if err == nil || err != io.EOF {
			return
		}

		_, err = br.fillBuffer()

		if err == io.EOF {
			err = nil
		} else if err != nil {
			return
		}

	case boundaryReaderStatePartialBoundaryInBuffer:
		if len(p) > br.remainingContent {
			p = p[:br.remainingContent+1]
		}

		n, err = br.buffer.Read(p)

		if err != nil {
			if err == io.EOF {
				err = errors.Errorf("unexpected EOF")
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		br.remainingContent -= n

		if br.remainingContent <= 0 {
			_, err = br.fillBuffer()

			if err != nil && err != io.EOF {
				return
			}
		}

	case boundaryReaderStateCompleteBoundaryInBuffer:
		if len(p) > br.remainingContent {
			p = p[:br.remainingContent+1]
		}

		n, err = br.buffer.Read(p)
		br.remainingContent -= n

		if err == io.EOF {
			_, err = br.fillBuffer()
		}

		if err == io.EOF {
			err = nil
		} else if err != nil {
			panic(fmt.Sprintf("invalid state: %q", err))
		}

		if br.remainingContent <= 0 {
			br.setState(boundaryReaderStateNeedsBoundary)
		}
	}

	return
}
