package gattungen

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

func init() {
	collections_value.RegisterGobValue[gattung.Gattung](nil)
}

type (
	Set        = schnittstellen.SetLike[gattung.Gattung]
	MutableSet = schnittstellen.MutableSetLike[gattung.Gattung]
)

func MakeSet(gs ...gattung.Gattung) Set {
	return collections_value.MakeValueSet[gattung.Gattung](nil, gs...)
}

func MakeMutableSet(gs ...gattung.Gattung) MutableSet {
	return collections_value.MakeMutableValueSet[gattung.Gattung](
		nil,
		gs...,
	)
}

func GattungFromString(v string) (s Set, err error) {
	parts := strings.Split(v, ",")
	partsGattung := make([]gattung.Gattung, len(parts))

	for i, v := range parts {
		if err = partsGattung[i].Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	s = MakeSet(partsGattung...)

	return
}

func GattungOrUnknownFromString(v string) (s Set) {
	parts := strings.Split(v, ",")
	partsGattung := make([]gattung.Gattung, len(parts))

	for i, v := range parts {
		if err := partsGattung[i].Set(v); err != nil {
			partsGattung[i] = gattung.Unknown
		}
	}

	s = MakeSet(partsGattung...)

	return
}
