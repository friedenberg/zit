package typ

import (
	"github.com/friedenberg/zit/src/alfa/etikett_rule"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/delta/kennung"
)

func Default() (t *Objekte, k *kennung.Typ) {
	k1 := kennung.MustTyp("md")
	k = &k1

	t = &Objekte{
		Akte: Akte{
			InlineAkte:     true,
			FileExtension:  "md",
			Formatters:     make(map[string]script_config.ScriptConfigWithUTI),
			Actions:        make(map[string]script_config.ScriptConfig),
			EtikettenRules: make(map[string]etikett_rule.Rule),
		},
	}

	return
}
