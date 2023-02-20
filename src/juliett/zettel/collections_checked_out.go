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

func ToSliceZettelsExternal(s SetCheckedOut) (out []External) {
	return collections.DerivedValues[CheckedOut, External](
		s,
		func(z CheckedOut) External {
			return z.External
		},
	)
}

func ToSliceFilesZettelen(s SetCheckedOut) (out []string) {
	return collections.DerivedValues[CheckedOut, string](
		s,
		func(z CheckedOut) string {
			return z.External.GetObjekteFD().Path
		},
	)
}

func ToSliceFilesAkten(s SetCheckedOut) (out []string) {
	return collections.DerivedValues[CheckedOut, string](
		s,
		func(z CheckedOut) string {
			return z.External.GetAkteFD().Path
		},
	)
}
