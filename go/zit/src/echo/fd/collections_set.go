package fd

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
)

func init() {
	collections_value.RegisterGobValue[*FD](nil)
}

type (
	Set        = interfaces.SetLike[*FD]
	MutableSet = interfaces.MutableSetLike[*FD]
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
	in interfaces.SetLike[T],
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
