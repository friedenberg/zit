package kennung

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
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
	return collections_value.MakeValueSet[FD](
		nil,
		ts...,
	)
}

func MakeMutableFDSet(ts ...FD) MutableFDSet {
	return collections_value.MakeMutableValueSet[FD](
		nil,
		ts...,
	)
}

func FDSetAddPairs[T FDPairGetter](
	in schnittstellen.SetLike[T],
	out schnittstellen.MutableSetLike[FD],
) (err error) {
	return in.Each(
		func(e T) (err error) {
			ofd := e.GetObjekteFD()

			if !ofd.IsEmpty() {
				if err = out.Add(ofd); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			ofd = e.GetAkteFD()

			if !ofd.IsEmpty() {
				if err = out.Add(ofd); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
	)
}

type FDKeyerSha struct{}

func (_ FDKeyerSha) GetKey(fd FD) string {
	return fd.Sha.String()
}
