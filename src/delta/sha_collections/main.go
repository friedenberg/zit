package sha_collections

import (
	"crypto/sha256"
	"io"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
)

type (
	Set        = schnittstellen.SetLike[sha.Sha]
	MutableSet = schnittstellen.MutableSetLike[sha.Sha]
)

func init() {
	collections.RegisterGob[sha.Sha]()
}

func MakeMutableSet(es ...sha.Sha) (s MutableSet) {
	return collections.MakeMutableSetStringer(es...)
}

func MakeMutableSetStrings(vs ...string) (s MutableSet, err error) {
	f := collections_ptr.MakeFlagCommas[sha.Sha, *sha.Sha](
		collections_ptr.SetterPolicyReset,
	)

	err = f.SetMany(vs...)
	s = f.CloneMutableSetLike()

	return
}

func ShaFromSet(s schnittstellen.SetLike[sha.Sha]) sha.Sha {
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
