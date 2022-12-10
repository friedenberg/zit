package typ

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/typ_toml"
)

func Default() (t *Objekte, k *kennung.Typ) {
	k1 := kennung.MustTyp("md")
	k = &k1

	t = &Objekte{
		Akte: typ_toml.Typ{
			InlineAkte:    true,
			FileExtension: "md",
		},
	}

	return
}
