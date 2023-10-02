package objekte_collections

import (
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func ToSliceFilesAkten(
	s sku.CheckedOutSet,
) (out []string, err error) {
	return iter.DerivedValuesPtr[sku.CheckedOut, *sku.CheckedOut, string](
		s,
		func(z *sku.CheckedOut) (e string, err error) {
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
	s sku.CheckedOutSet,
) (out []string, err error) {
	return iter.DerivedValuesPtr[sku.CheckedOut, *sku.CheckedOut, string](
		s,
		func(z *sku.CheckedOut) (e string, err error) {
			e = z.External.GetObjekteFD().Path

			if e == "" {
				err = collections.MakeErrStopIteration()
				return
			}

			return
		},
	)
}
