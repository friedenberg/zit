package id_fmts

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_ptr"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type etikettenReader struct{}

func MakeEtikettenReader() (f *etikettenReader) {
	f = &etikettenReader{}

	return
}

func (f *etikettenReader) ReadStringFormat(
	rb *catgut.RingBuffer,
	k ids.TagMutableSet,
) (n int64, err error) {
	flag := collections_ptr.MakeFlagCommasFromExisting(
		collections_ptr.SetterPolicyAppend,
		k,
	)

	var readable catgut.Slice

	if readable, err = rb.PeekUptoAndIncluding('\n'); err != nil && err != io.EOF {
		errors.Wrap(err)
		return
	}

	if readable.Len() == 1 {
		return
	}

	if err = flag.Set(readable.String()); err != nil {
		errors.Wrap(err)
		return
	}

	n = int64(readable.Len())
	rb.AdvanceRead(readable.Len())

	return
}
