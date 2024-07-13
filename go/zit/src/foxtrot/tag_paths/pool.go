package tag_paths

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var p interfaces.Pool[Path, *Path]

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

func GetPool() interfaces.Pool[Path, *Path] {
	return p
}
