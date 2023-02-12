package sha_collections

import (
	"crypto/sha256"
	"io"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type (
	Set        = collections.ValueSet[sha.Sha, *sha.Sha]
	MutableSet = collections.MutableValueSet[sha.Sha, *sha.Sha]
)

func MakeMutableSet(es ...sha.Sha) (s MutableSet) {
	return MutableSet(collections.MakeMutableValueSet(es...))
}

func MakeMutableSetStrings(vs ...string) (s MutableSet, err error) {
	var s1 collections.MutableValueSet[sha.Sha, *sha.Sha]

	if s1, err = collections.MakeMutableValueSetStrings[sha.Sha, *sha.Sha](vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	s = MutableSet(s1)

	return
}

func ShaFromSet(s collections.SetLike[sha.Sha]) sha.Sha {
	hash := sha256.New()

	elements := make([]sha.Sha, 0, s.Len())

	s.Each(
		func(s sha.Sha) (err error) {
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

	return sha.FromHash(hash)
}
