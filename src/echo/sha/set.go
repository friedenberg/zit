package sha

import (
	"crypto/sha256"
	"io"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/collections"
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

func ShaFromSet(s collections.SetLike[Sha]) Sha {
	hash := sha256.New()

	elements := make([]Sha, 0, s.Len())

	s.Each(
		func(s Sha) (err error) {
			elements = append(elements, s)
			return
		},
	)

	sort.Slice(
		elements,
		func(i, j int) bool { return elements[i].String() < elements[j].String() },
	)

	for _, e := range elements {
		if _, err := io.WriteString(hash, e.String()); err != nil {
			errors.PanicIfError(errors.Wrap(err))
		}
	}

	return FromHash(hash)
}
