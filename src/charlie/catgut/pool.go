package catgut

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/pool"
)

var p schnittstellen.Pool[String, *String]

func init() {
	p = pool.MakePool[String, *String](
		nil,
		func(v *String) {
			v.Reset()
		},
	)
}

func GetPool() schnittstellen.Pool[String, *String] {
	return p
}
