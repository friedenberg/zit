package flag

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/flag_policy"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
)

type Flag[T interface {
	schnittstellen.StringerSetter
	schnittstellen.Resetter
}] struct {
	flag_policy.FlagPolicy
	Embedded T
}

func (f Flag[T]) Set(v string) (err error) {
	if f.FlagPolicy == flag_policy.FlagPolicyReset {
		f.Embedded.Reset()
	}

	if err = f.Embedded.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f Flag[T]) String() string {
	return f.Embedded.String()
}
