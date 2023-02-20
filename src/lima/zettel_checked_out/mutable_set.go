package zettel_checked_out

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type (
	Set        = schnittstellen.Set[Zettel]
	MutableSet = schnittstellen.MutableSet[Zettel]
)

func MakeMutableSetUnique(c int) MutableSet {
	return collections.MakeMutableSet(
		func(sz Zettel) string {
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

func MakeMutableSetHinweisZettel(c int) MutableSet {
	return collections.MakeMutableSet(
		func(sz Zettel) string {
			return collections.MakeKey(sz.Internal.Sku.Kennung)
		},
	)
}

func ToSliceZettelsExternal(s Set) (out []zettel.External) {
	return collections.DerivedValues[Zettel, zettel.External](
		s,
		func(z Zettel) zettel.External {
			return z.External
		},
	)
}

func ToSliceFilesZettelen(s Set) (out []string) {
	return collections.DerivedValues[Zettel, string](
		s,
		func(z Zettel) string {
			return z.External.GetObjekteFD().Path
		},
	)
}

func ToSliceFilesAkten(s Set) (out []string) {
	return collections.DerivedValues[Zettel, string](
		s,
		func(z Zettel) string {
			return z.External.GetAkteFD().Path
		},
	)
}
