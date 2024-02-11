package standort

import (
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
)

func (s Standort) ObjekteReaderFactory(
	g schnittstellen.GattungGetter,
) schnittstellen.ObjekteReaderFactory {
	return s.ObjekteReaderWriterFactory(g)
}

func (s Standort) ObjekteReaderWriterFactory(
	g schnittstellen.GattungGetter,
) schnittstellen.ObjekteIOFactory {
	return schnittstellen.MakeBespokeObjekteReadWriterFactory(
		schnittstellen.MakeBespokeObjekteReadFactory(
			func(sh schnittstellen.ShaGetter) (schnittstellen.ShaReadCloser, error) {
				return s.objekteReader(g, sh)
			},
		),
		schnittstellen.MakeBespokeObjekteWriteFactory(
			func() (schnittstellen.ShaWriteCloser, error) {
				return s.objekteWriter(g)
			},
		),
	)
}
