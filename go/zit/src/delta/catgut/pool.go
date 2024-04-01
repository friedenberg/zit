package catgut

import (
	"sync"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/pool"
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
