package sha

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
)

type Set = collections.ValueSet[Sha, *Sha]
type MutableSet = collections.MutableValueSet[Sha, *Sha]

func MakeMutableSet(es ...Sha) (s MutableSet) {
	return MutableSet(collections.MakeMutableValueSet(es...))
}

func MakeMutableSetStrings(vs ...string) (s MutableSet, err error) {
	var s1 collections.MutableValueSet[Sha, *Sha]

	if s1, err = collections.MakeMutableValueSetStrings[Sha, *Sha](vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	s = MutableSet(s1)

	return
}
