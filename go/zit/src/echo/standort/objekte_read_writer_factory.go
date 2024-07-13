package standort

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func (s Standort) ObjekteReaderFactory(
	g interfaces.GattungGetter,
) interfaces.ObjekteReaderFactory {
	return s.ObjekteReaderWriterFactory(g)
}

func (s Standort) ObjekteReaderWriterFactory(
	g interfaces.GattungGetter,
) interfaces.ObjekteIOFactory {
	return interfaces.MakeBespokeObjekteReadWriterFactory(
		interfaces.MakeBespokeObjekteReadFactory(
			func(sh interfaces.ShaGetter) (interfaces.ShaReadCloser, error) {
				return s.objekteReader(g, sh)
			},
		),
		interfaces.MakeBespokeObjekteWriteFactory(
			func() (interfaces.ShaWriteCloser, error) {
				return s.objekteWriter(g)
			},
		),
	)
}
