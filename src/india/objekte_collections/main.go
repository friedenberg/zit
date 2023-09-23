package objekte_collections

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

func ToSliceFilesAkten(
	s schnittstellen.SetLike[*objekte.CheckedOut],
) (out []string, err error) {
	return iter.DerivedValues[*objekte.CheckedOut, string](
		s,
		func(z *objekte.CheckedOut) (e string, err error) {
			e = z.External.GetAkteFD().Path

			if e == "" {
				err = collections.MakeErrStopIteration()
				return
			}

			return
		},
	)
}

func ToSliceFilesZettelen(
	s schnittstellen.SetLike[*objekte.CheckedOut],
) (out []string, err error) {
	return iter.DerivedValues[*objekte.CheckedOut, string](
		s,
		func(z *objekte.CheckedOut) (e string, err error) {
			e = z.External.GetObjekteFD().Path

			if e == "" {
				err = collections.MakeErrStopIteration()
				return
			}

			return
		},
	)
}
