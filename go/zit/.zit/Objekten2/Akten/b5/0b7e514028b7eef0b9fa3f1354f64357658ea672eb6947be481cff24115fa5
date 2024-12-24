package sha_probe_index

import (
	"bytes"
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

const RowSize = sha.ByteSize + sha.ByteSize

type row struct {
	left  sha.Sha
	right sha.Sha
}

func (r *row) IsEmpty() bool {
	return r.right.IsNull() && r.left.IsNull()
}

func (r *row) String() string {
	return fmt.Sprintf(
		"%s %s",
		r.left.GetShaString()[:4],
		r.right.GetShaString()[:4],
	)
}

func (current *row) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int64

	n1, err = current.left.ReadFrom(r)
	n += n1

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n1, err = current.right.ReadFrom(r)
	n += int64(n1)

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	if n != RowSize {
		err = errors.Errorf("expected to read %d but read %d", RowSize, n)
		return
	}

	return
}

func (r *row) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int
	var n2 int64

	n, err = r.left.WriteTo(w)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = r.right.WriteTo(w)
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if n != RowSize {
		err = errors.Errorf("expected to write %d but read %d", RowSize, n)
		return
	}

	return
}

type rowEqualerComplete struct{}

func (rowEqualerComplete) Equals(a, b *row) bool {
	return a.left.Equals(&b.left) && a.right.Equals(&b.right)
}

type rowEqualerShaOnly struct{}

func (rowEqualerShaOnly) Equals(a, b *row) bool {
	return a.left.Equals(&b.left)
}

type rowResetter struct{}

func (rowResetter) Reset(a *row) {
	a.left.Reset()
	a.right.Reset()
}

func (rowResetter) ResetWith(a, b *row) {
	a.left.ResetWith(&b.left)
	a.right.ResetWith(&b.right)
}

var RowLessor rowLessor

type rowLessor struct{}

func (rowLessor) Less(a, b *row) bool {
	cmp := bytes.Compare(a.left.GetShaBytes(), b.left.GetShaBytes())

	if cmp != 0 {
		return cmp == -1
	}

	cmp = bytes.Compare(a.right.GetShaBytes(), b.right.GetShaBytes())

	return cmp == -1
}
