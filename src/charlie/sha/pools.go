package sha

import (
	"bytes"
	"crypto/sha256"
	"hash"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/pool"
)

var shaPool schnittstellen.PoolValue[hash.Hash]

func init() {
	shaPool = pool.MakePoolValue[hash.Hash](
		func() hash.Hash {
			return sha256.New()
		},
		func(h hash.Hash) {
			h.Reset()
		},
	)
}

var Resetter resetter

type resetter struct{}

func (resetter) Reset(s *Sha) {
	s.Reset()
}

func (resetter) ResetWith(a, b *Sha) {
	a.ResetWith(b)
}

var Lessor lessor

type lessor struct{}

func (lessor) Less(a, b *Sha) bool {
	return bytes.Compare(a.GetShaBytes(), b.GetShaBytes()) == -1
}

var Equaler equaler

type equaler struct{}

func (equaler) Equals(a, b *Sha) bool {
	return bytes.Equal(a.GetShaBytes(), b.GetShaBytes())
}
