package fs_home

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func (s Home) ObjekteReaderFactory(
	g interfaces.GenreGetter,
) interfaces.ObjectReaderFactory {
	return s.ObjekteReaderWriterFactory(g)
}

func (s Home) ObjekteReaderWriterFactory(
	g interfaces.GenreGetter,
) interfaces.ObjectIOFactory {
	return interfaces.MakeBespokeObjectReadWriterFactory(
		interfaces.MakeBespokeObjectReadFactory(
			func(sh interfaces.ShaGetter) (interfaces.ShaReadCloser, error) {
				return s.objectReader(g, sh)
			},
		),
		interfaces.MakeBespokeObjectWriteFactory(
			func() (interfaces.ShaWriteCloser, error) {
				return s.objectWriter(g)
			},
		),
	)
}
