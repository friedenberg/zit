package sha

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/proto_objekte"
)

type MutableSet = proto_objekte.MutableSet[Sha, *Sha]

func MakeMutableSet(es ...Sha) (s MutableSet) {
	return MutableSet(proto_objekte.MakeMutableSet(es...))
}

func MakeMutableSetStrings(vs ...string) (s MutableSet, err error) {
	var s1 proto_objekte.MutableSet[Sha, *Sha]

	if s1, err = proto_objekte.MakeMutableSetStrings[Sha, *Sha](vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	s = MutableSet(s1)

	return
}
