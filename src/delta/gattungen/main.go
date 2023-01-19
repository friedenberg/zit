package gattungen

import (
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
