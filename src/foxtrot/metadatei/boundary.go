package metadatei

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/catgut"
)

const (
	Boundary = "---"
)

var (
	BoundaryStringValue values.String
	boundaryBytes       = []byte(Boundary)
)

func init() {
	BoundaryStringValue = values.MakeString(Boundary)
}

var errBoundaryInvalid = errors.New("boundary invalid")

func ReadBoundary(r *catgut.RingBuffer) (n int, err error) {
	var readable catgut.Slice

	readable, err = r.PeekUpto('\n')

	if err != nil && err != io.EOF {
		return
	} else if r.PeekReadable().Len() == 0 && err == io.EOF {
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
