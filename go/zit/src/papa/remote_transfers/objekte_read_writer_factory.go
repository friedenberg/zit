package remote_transfers

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
)

func (s client) ObjekteReaderFactory(
	g schnittstellen.GattungGetter,
) schnittstellen.ObjekteReaderFactory {
	return schnittstellen.MakeBespokeObjekteReadFactory(
		func(sh schnittstellen.ShaGetter) (schnittstellen.ShaReadCloser, error) {
			return s.ObjekteReader(g, sh)
		},
	)
}
