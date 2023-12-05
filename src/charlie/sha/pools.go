package sha

import (
	"crypto/sha256"
	"hash"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/pool"
)

var (
	shaPool schnittstellen.PoolValue[hash.Hash]
)

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
