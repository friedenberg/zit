package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_value"
)

func init() {
	collections_value.RegisterGobValue[FD](nil)
}

type (
	FDSet        = schnittstellen.SetLike[FD]
	MutableFDSet = schnittstellen.MutableSetLike[FD]
)

func MakeFDSet(ts ...FD) FDSet {
	return collections.MakeSet[FD](
		(FD).String,
		ts...,
	)
}

func MakeMutableFDSet(ts ...FD) MutableFDSet {
	return collections.MakeMutableSet[FD](
		(FD).String,
		ts...,
	)
}

func FDSetAddPairs[T FDPairGetter](
	in schnittstellen.SetLike[T],
	out schnittstellen.MutableSetLike[FD],
) (err error) {
	return in.Each(
		func(e T) (err error) {
			out.Add(e.GetObjekteFD())
			out.Add(e.GetAkteFD())
			return
		},
	)
}

func FDSetContainsPair(s FDSet, maybeFDs Matchable) (ok bool) {
	var fdGetter FDPairGetter

	if fdGetter, ok = maybeFDs.(FDPairGetter); !ok {
		return
	}

	objekte := fdGetter.GetObjekteFD()

	if ok = s.Contains(objekte); ok {
		return
	}

	akte := fdGetter.GetAkteFD()

	if ok = s.Contains(akte); ok {
		return
	}

	return
}
