package kennung_fmt

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type etikettenReader struct{}

func MakeEtikettenReader() (f *etikettenReader) {
	f = &etikettenReader{}

	return
}

func (f *etikettenReader) ReadStringFormat(
	rb *catgut.RingBuffer,
	k kennung.EtikettMutableSet,
) (n int64, err error) {
	flag := collections_ptr.MakeFlagCommasFromExisting(
		collections_ptr.SetterPolicyAppend,
		k,
	)

	readable, err := rb.PeekUpto('\n')

	if err = flag.Set(readable.String()); err != nil {
		errors.Wrap(err)
		return
	}

	return
}
