package remote_transfers

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func (s client) ObjectReaderFactory(
	g interfaces.GenreGetter,
) interfaces.ObjectReaderFactory {
	return interfaces.MakeBespokeObjectReadFactory(
		func(sh interfaces.ShaGetter) (interfaces.ShaReadCloser, error) {
			return s.ObjectReader(g, sh)
		},
	)
}
