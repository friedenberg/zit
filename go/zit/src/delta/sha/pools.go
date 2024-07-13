package sha

import (
	"bytes"
	"crypto/sha256"
	"hash"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var (
	hash256Pool interfaces.PoolValue[hash.Hash]
	shaPool     interfaces.Pool[Sha, *Sha]
)

func init() {
	hash256Pool = pool.MakeValue[hash.Hash](
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

func GetPool() interfaces.Pool[Sha, *Sha] {
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
