package objekte_collections

import (
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

func ToSliceFilesZettelen(
	s sku.CheckedOutSet,
) (out []string, err error) {
	return iter.DerivedValues(
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
