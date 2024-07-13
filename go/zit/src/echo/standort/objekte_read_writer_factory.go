package standort

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func (s Standort) ObjekteReaderFactory(
	g interfaces.GenreGetter,
) interfaces.ObjectReaderFactory {
	return s.ObjekteReaderWriterFactory(g)
}

func (s Standort) ObjekteReaderWriterFactory(
	g interfaces.GenreGetter,
) interfaces.ObjectIOFactory {
	return interfaces.MakeBespokeObjectReadWriterFactory(
		interfaces.MakeBespokeObjectReadFactory(
			func(sh interfaces.ShaGetter) (interfaces.ShaReadCloser, error) {
				return s.objekteReader(g, sh)
			},
		),
		interfaces.MakeBespokeObjectWriteFactory(
			func() (interfaces.ShaWriteCloser, error) {
				return s.objekteWriter(g)
			},
		),
	)
}
