package ohio

import (
	"bufio"
	"bytes"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/log"
)

type BoundaryReader interface {
	io.Reader
	ReadBoundary() (int, error)
}

type content []byte

type boundaryReader struct {
	br       *bufio.Reader
	boundary []byte
	// make circular
	head                                []byte
	tail                                []byte
	needsBoundaryAfterCompletedHeadRead bool
}

func MakeBoundaryReader(r io.Reader, boundary string) BoundaryReader {
	return &boundaryReader{
		br:                                  bufio.NewReader(r),
		boundary:                            []byte(boundary),
		needsBoundaryAfterCompletedHeadRead: true,
	}
}

func (br *boundaryReader) ReadBoundary() (n int, err error) {
	if !br.needsBoundary() {
		err = errors.Errorf("next read should be content, not boundary")
		return
	}

	log.Log().Printf("read boundary")
	br.needsBoundaryAfterCompletedHeadRead = false

	return
}

func (br *boundaryReader) resectHead() (found bool) {
	if len(br.head) > 0 {
		panic("resecting non-empty head")
	}

	var head, tail []byte

	head, tail, found = bytes.Cut(br.tail, br.boundary)

	br.head = append(br.head, head...)
	br.tail = tail

	return
}

func (br *boundaryReader) readFromBufferIfNecessary(
	p []byte,
) (n int) {
	if len(br.head) == 0 {
		return
	}

	// copy up to max length of p
	lenB := len(br.head)
	lenP := len(p)
	nToCopy := lenB

	if lenB > lenP {
		nToCopy = lenP
	}

	n = copy(p, br.head[:nToCopy])
	br.head = br.head[nToCopy:]

	if len(br.head) == 0 {
		br.needsBoundaryAfterCompletedHeadRead = true
		br.resectHead()
	}

	return
}

func (br *boundaryReader) needsBoundary() bool {
	if len(br.head) > 0 {
		return false
	}

	if !br.needsBoundaryAfterCompletedHeadRead {
		return false
	}

	return true
}

func (br *boundaryReader) Read(p []byte) (n int, err error) {
	if br.needsBoundary() {
		err = io.EOF
		// err = errors.Errorf("next read should be boundary, not content")
		return
	}

	n = br.readFromBufferIfNecessary(p)

	if n > 0 {
		return
	}

	n, err = br.br.Read(p)

	found := br.fillBuffer(p)

	if found {
		n = br.readFromBufferIfNecessary(p)
		return
	}

	return
}

func (br *boundaryReader) fillBuffer(p []byte) (found bool) {
	if len(br.head) > 0 {
		br.tail = append(br.tail, p...)
		return
	}

	br.tail = append(br.tail, p...)
	found = br.resectHead()

	return
}
