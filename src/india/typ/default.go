package typ

import (
	"github.com/friedenberg/zit/src/alfa/etikett_rule"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
)

func Default() (t *Objekte, k *kennung.Typ) {
	k1 := kennung.MustTyp("md")
	k = &k1

	t = &Objekte{
		Akte: Akte{
			InlineAkte:     true,
			FileExtension:  "md",
			Actions:        make(map[string]Action),
			EtikettenRules: make(map[string]etikett_rule.Rule),
		},
	}

	return
}
