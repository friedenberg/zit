package sha

import (
	"bytes"
	"crypto/sha256"
	"hash"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/pool"
)

var (
	hash256Pool schnittstellen.PoolValue[hash.Hash]
	shaPool     schnittstellen.Pool[Sha, *Sha]
)

func init() {
	hash256Pool = pool.MakePoolValue[hash.Hash](
		func() hash.Hash {
			return sha256.New()
		},
		func(h hash.Hash) {
			h.Reset()
		},
	)
	shaPool = pool.MakePool[Sha, *Sha](
		nil,
		func(sh *Sha) {
			sh.Reset()
		},
	)
}

func GetPool() schnittstellen.Pool[Sha, *Sha] {
	return shaPool
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
