package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
)

func ToSliceFilesZettelen(
	s CheckedOutFSSet,
) (out []string, err error) {
	return iter.DerivedValues(
		s,
		func(z *CheckedOutFS) (e string, err error) {
			e = z.External.GetObjekteFD().GetPath()

			if e == "" {
				err = collections.MakeErrStopIteration()
				return
			}

			return
		},
	)
}
