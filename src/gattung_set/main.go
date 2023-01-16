package gattung_set

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
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
