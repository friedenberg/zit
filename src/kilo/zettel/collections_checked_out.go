package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/kilo/checked_out"
)

type (
	SetCheckedOut        = schnittstellen.SetLike[checked_out.Zettel]
	MutableSetCheckedOut = schnittstellen.MutableSetLike[checked_out.Zettel]
)

type CheckedOutUniqueKeyer struct{}

func (k CheckedOutUniqueKeyer) GetKey(sz checked_out.Zettel) string {
	return collections.MakeKey(
		sz.Internal.GetTai(),
		sz.Internal.GetKennung(),
		sz.Internal.ObjekteSha,
	)
}

func MakeMutableSetCheckedOutUnique(c int) MutableSetCheckedOut {
	return collections_value.MakeMutableValueSet[checked_out.Zettel](
		CheckedOutUniqueKeyer{},
	)
}

type CheckedOutHinweisKeyer struct{}

func (k CheckedOutHinweisKeyer) GetKey(sz checked_out.Zettel) string {
	return collections.MakeKey(
		sz.Internal.GetKennung(),
	)
}

func MakeMutableSetCheckedOutHinweisZettel(c int) MutableSetCheckedOut {
	return collections_value.MakeMutableValueSet[checked_out.Zettel](
		CheckedOutHinweisKeyer{},
	)
}

func ToSliceFilesZettelen(s SetCheckedOut) (out []string, err error) {
	return iter.DerivedValues[checked_out.Zettel, string](
		s,
		func(z checked_out.Zettel) (e string, err error) {
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
	return iter.DerivedValues[checked_out.Zettel, string](
		s,
		func(z checked_out.Zettel) (e string, err error) {
			e = z.External.GetAkteFD().Path

			if e == "" {
				err = collections.MakeErrStopIteration()
				return
			}

			return
		},
	)
}
