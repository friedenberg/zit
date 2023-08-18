package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_value"
)

type (
	SetCheckedOut        = schnittstellen.SetLike[CheckedOut]
	MutableSetCheckedOut = schnittstellen.MutableSetLike[CheckedOut]
)

type CheckedOutUniqueKeyer struct{}

func (k CheckedOutUniqueKeyer) GetKey(sz CheckedOut) string {
	return collections.MakeKey(
		sz.Internal.Kopf,
		sz.Internal.GetTai(),
		sz.Internal.GetKennung(),
		sz.Internal.ObjekteSha,
	)
}

func MakeMutableSetCheckedOutUnique(c int) MutableSetCheckedOut {
	return collections_value.MakeMutableValueSet[CheckedOut](
		CheckedOutUniqueKeyer{},
	)
}

type CheckedOutHinweisKeyer struct{}

func (k CheckedOutHinweisKeyer) GetKey(sz CheckedOut) string {
	return collections.MakeKey(
		sz.Internal.GetKennung(),
	)
}

func MakeMutableSetCheckedOutHinweisZettel(c int) MutableSetCheckedOut {
	return collections_value.MakeMutableValueSet[CheckedOut](
		CheckedOutHinweisKeyer{},
	)
}

func ToSliceFilesZettelen(s SetCheckedOut) (out []string, err error) {
	return iter.DerivedValues[CheckedOut, string](
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
	return iter.DerivedValues[CheckedOut, string](
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
