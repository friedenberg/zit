package sha

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	collections "github.com/friedenberg/zit/src/bravo/collections"
)

type Set = collections.Set[Sha, *Sha]
type MutableSet = collections.MutableSet[Sha, *Sha]

func MakeMutableSet(es ...Sha) (s MutableSet) {
	return MutableSet(collections.MakeMutableSet(es...))
}

func MakeMutableSetStrings(vs ...string) (s MutableSet, err error) {
	var s1 collections.MutableSet[Sha, *Sha]

	if s1, err = collections.MakeMutableSetStrings[Sha, *Sha](vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	s = MutableSet(s1)

	return
}
