package remote_transfers

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func (s client) ObjekteReaderFactory(
	g interfaces.GattungGetter,
) interfaces.ObjekteReaderFactory {
	return interfaces.MakeBespokeObjekteReadFactory(
		func(sh interfaces.ShaGetter) (interfaces.ShaReadCloser, error) {
			return s.ObjekteReader(g, sh)
		},
	)
}
