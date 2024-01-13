package catgut

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/pool"
)

var (
	p     schnittstellen.Pool[String, *String]
	ponce sync.Once
)

func init() {
}

func GetPool() schnittstellen.Pool[String, *String] {
	ponce.Do(
		func() {
			p = pool.MakePool[String, *String](
				nil,
				func(v *String) {
					v.Reset()
				},
			)
		},
	)

	return p
}
