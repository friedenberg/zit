package dir_layout

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func (s DirLayout) ObjectReaderWriterFactory(
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
