package gattungen

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
)

func init() {
	collections.RegisterGob[gattung.Gattung]()
}

type (
	Set        = schnittstellen.Set[gattung.Gattung]
	MutableSet = schnittstellen.MutableSet[gattung.Gattung]
)

func MakeSet(gs ...gattung.Gattung) Set {
	return collections.MakeSet[gattung.Gattung](
		func(g gattung.Gattung) string {
			return g.String()
		},
		gs...,
	)
}

func MakeMutableSet(gs ...gattung.Gattung) MutableSet {
	return collections.MakeMutableSet[gattung.Gattung](
		func(g gattung.Gattung) string {
			return g.String()
		},
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
