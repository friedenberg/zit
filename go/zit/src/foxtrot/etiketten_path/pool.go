package etiketten_path

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var p schnittstellen.Pool[Path, *Path]

func init() {
	p = pool.MakePool(
		func() *Path {
			return &Path{}
		},
		func(p *Path) {
			for _, s := range *p {
				s.Reset()
			}

			*p = (*p)[:0]
		},
	)
}

func GetPool() schnittstellen.Pool[Path, *Path] {
	return p
}
