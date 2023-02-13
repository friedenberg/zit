package collections

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type SetterPolicy int

const (
	SetterPolicyAppend = SetterPolicy(iota)
	SetterPolicyReset
)

type MutableValueSet2[
	E schnittstellen.Value,
	EPtr schnittstellen.ValuePtr[E],
] struct {
	schnittstellen.MutableSetLike[E]
	SetterPolicy
}

func (vs MutableValueSet2[E, EPtr]) Set(v string) (err error) {
	switch vs.SetterPolicy {
	case SetterPolicyReset:
		vs.MutableSetLike.Reset(nil)
	}

	els := strings.Split(v, ",")

	for _, e := range els {
		e = strings.TrimSpace(e)

		if err = AddString[E, EPtr](vs.MutableSetLike, e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (v MutableValueSet2[E, EPtr]) String() string {
	return String[E](v.MutableSetLike)
}
