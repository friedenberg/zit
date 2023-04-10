package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type (
	SetCheckedOut        = schnittstellen.Set[CheckedOut]
	MutableSetCheckedOut = schnittstellen.MutableSet[CheckedOut]
)

func MakeMutableSetCheckedOutUnique(c int) MutableSetCheckedOut {
	return collections.MakeMutableSet(
		func(sz CheckedOut) string {
			return collections.MakeKey(
				sz.Internal.Sku.Kopf,
				sz.Internal.Sku.Mutter[0],
				sz.Internal.Sku.Mutter[1],
				sz.Internal.Sku.Schwanz,
				sz.Internal.Sku.Kennung,
				sz.Internal.Sku.ObjekteSha,
			)
		},
	)
}

func MakeMutableSetCheckedOutHinweisZettel(c int) MutableSetCheckedOut {
	return collections.MakeMutableSet(
		func(sz CheckedOut) string {
			return collections.MakeKey(sz.Internal.Sku.Kennung)
		},
	)
}

func ToSliceFilesZettelen(s SetCheckedOut) (out []string, err error) {
	return collections.DerivedValues[CheckedOut, string](
		s,
		func(z CheckedOut) (e string, err error) {
			e = z.External.GetObjekteFD().Path

			if e == "" {
				err = collections.MakeErrStopIteration()
				return
			}

			return
		},
	)
}

func ToSliceFilesAkten(s SetCheckedOut) (out []string, err error) {
	return collections.DerivedValues[CheckedOut, string](
		s,
		func(z CheckedOut) (e string, err error) {
			e = z.External.GetAkteFD().Path

			if e == "" {
				err = collections.MakeErrStopIteration()
				return
			}

			return
		},
	)
}
