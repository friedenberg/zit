package triple_hyphen_io

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
)

const (
	Boundary = "---"
)

var BoundaryStringValue catgut.String

func init() {
	errors.PanicIfError(BoundaryStringValue.Set(Boundary))
}

var errBoundaryInvalid = errors.New("boundary invalid")

func ReadBoundary(r *catgut.RingBuffer) (n int, err error) {
	var readable catgut.Slice

	readable, err = r.PeekUpto('\n')

	if errors.IsNotNilAndNotEOF(err) {
		return
	} else if readable.Len() == 0 && err == io.EOF {
		return
	}

	if readable.Len() != BoundaryStringValue.Len() {
		err = errors.Wrapf(
			errBoundaryInvalid,
			"Expected: %d, Actual: %d. Readable: %q",
			BoundaryStringValue.Len(),
			readable.Len(),
			readable,
		)

		return
	}

	if !readable.Equal(BoundaryStringValue.Bytes()) {
		err = errors.Wrapf(
			errBoundaryInvalid,
			"Expected: %q, Actual: %q",
			BoundaryStringValue.String(),
			readable.String(),
		)
		return
	}

	// boundary and newline
	n = BoundaryStringValue.Len() + 1
	r.AdvanceRead(n)

	if err == io.EOF {
		err = nil
	}

	return
}
