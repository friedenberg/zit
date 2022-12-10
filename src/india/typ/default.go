package typ

import (
	"github.com/friedenberg/zit/src/foxtrot/kennung"
)

func Default() (t *Objekte, k *kennung.Typ) {
	k1 := kennung.MustTyp("md")
	k = &k1

	t = &Objekte{
		Akte: Akte{
			InlineAkte:    true,
			FileExtension: "md",
		},
	}

	return
}
