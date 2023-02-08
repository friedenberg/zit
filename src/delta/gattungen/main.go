package gattungen

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type Set = collections.Set[gattung.Gattung]
type MutableSet = collections.MutableSet[gattung.Gattung]

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

// type Setter struct{
//   MutableSet
//   append bool
// }

// func (v Setter) String() string {
//   return v.MutableSet.String()
// }

// func (v *Setter) Set(v string) (err error) {
//   return
// }

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
