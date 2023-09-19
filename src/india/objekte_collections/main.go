package objekte_collections

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

func ToSliceFilesAkten(
	s schnittstellen.SetLike[objekte.CheckedOutLikePtr],
) (out []string, err error) {
	return iter.DerivedValues[objekte.CheckedOutLikePtr, string](
		s,
		func(z objekte.CheckedOutLikePtr) (e string, err error) {
			e = z.GetExternalLike().GetAkteFD().Path

			if e == "" {
				err = collections.MakeErrStopIteration()
				return
			}

			return
		},
	)
}

func ToSliceFilesZettelen(
	s schnittstellen.SetLike[objekte.CheckedOutLikePtr],
) (out []string, err error) {
	return iter.DerivedValues[objekte.CheckedOutLikePtr, string](
		s,
		func(z objekte.CheckedOutLikePtr) (e string, err error) {
			e = z.GetExternalLike().GetObjekteFD().Path

			if e == "" {
				err = collections.MakeErrStopIteration()
				return
			}

			return
		},
	)
}
