package metadatei

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/ohio_ring_buffer2"
)

var (
	BoundaryStringValue values.String
	boundaryBytes       = []byte{'-', '-', '-'}
)

const (
	Boundary = "---"
)

func init() {
	BoundaryStringValue = values.MakeString(Boundary)
}

var errBoundaryInvalid = errors.New("boundary invalid")

func ReadBoundary(r *ohio_ring_buffer2.RingBuffer) (n int, err error) {
	var (
		readable ohio_ring_buffer2.Slice
		ok       bool
	)

	readable, ok, err = r.PeekUpto('\n')

	if !ok {
		if r.PeekReadable().Len() == 0 {
			err = io.EOF
		} else {
			err = errBoundaryInvalid
		}

		return
	}

	if !readable.Equal(boundaryBytes) {
		err = errBoundaryInvalid
		return
	}

	// boundary and newline
	n = len(boundaryBytes) + 1
	r.AdvanceRead(n)

	if err == io.EOF {
		err = nil
	}

	return
}
