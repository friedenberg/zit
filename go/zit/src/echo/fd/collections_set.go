package fd

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_value"
)

func init() {
	collections_value.RegisterGobValue[*FD](nil)
}

type (
	Set        = schnittstellen.SetLike[*FD]
	MutableSet = schnittstellen.MutableSetLike[*FD]
)

func MakeSet(ts ...*FD) Set {
	return collections_value.MakeValueSet[*FD](
		nil,
		ts...,
	)
}

func MakeMutableSet(ts ...*FD) MutableSet {
	return collections_value.MakeMutableValueSet[*FD](
		nil,
		ts...,
	)
}

func MakeMutableSetSha() MutableSet {
	return collections_value.MakeMutableValueSet[*FD](
		KeyerSha{},
	)
}

func SetAddPairs[T FDPairGetter](
	in schnittstellen.SetLike[T],
	out MutableSet,
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

type KeyerSha struct{}

func (KeyerSha) GetKey(fd *FD) string {
	return fd.sha.String()
}
