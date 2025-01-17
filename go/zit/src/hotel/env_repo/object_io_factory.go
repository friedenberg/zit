package env_repo

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

// TODO remove this in favor of using an append-only locked log of inventory
// lists
func (s Env) ObjectReaderWriterFactory(
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
