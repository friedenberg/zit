package remote_transfers

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func (s client) ObjekteReaderFactory(
	g interfaces.GenreGetter,
) interfaces.ObjectReaderFactory {
	return interfaces.MakeBespokeObjectReadFactory(
		func(sh interfaces.ShaGetter) (interfaces.ShaReadCloser, error) {
			return s.ObjekteReader(g, sh)
		},
	)
}
