package objekte_collections

import (
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func ToSliceFilesAkten(
	s sku.CheckedOutSet,
) (out []string, err error) {
	return iter.DerivedValues[*sku.CheckedOut, string](
		s,
		func(z *sku.CheckedOut) (e string, err error) {
			e = z.External.GetAkteFD().GetPath()

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
	return iter.DerivedValues[*sku.CheckedOut, string](
		s,
		func(z *sku.CheckedOut) (e string, err error) {
			e = z.External.GetObjekteFD().GetPath()

			if e == "" {
				err = collections.MakeErrStopIteration()
				return
			}

			return
		},
	)
}
