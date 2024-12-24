package catgut

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var (
	p     interfaces.Pool[String, *String]
	ponce sync.Once
)

func init() {
}

func GetPool() interfaces.Pool[String, *String] {
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
