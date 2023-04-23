package typ

import (
	"github.com/friedenberg/zit/src/alfa/etikett_rule"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/delta/kennung"
)

func Default() (t Objekte, k kennung.Typ) {
	k = kennung.MustTyp("md")

	t = Objekte{
		Akte: Akte{
			InlineAkte:     true,
			FileExtension:  "md",
			Formatters:     make(map[string]script_config.WithOutputFormat),
			VimSyntaxType:  "markdown",
			Actions:        make(map[string]script_config.ScriptConfig),
			EtikettenRules: make(map[string]etikett_rule.Rule),
		},
	}

	return
}

func MakeObjekte() (t *Objekte) {
	t = &Objekte{
		Akte: Akte{
			Formatters:     make(map[string]script_config.WithOutputFormat),
			Actions:        make(map[string]script_config.ScriptConfig),
			EtikettenRules: make(map[string]etikett_rule.Rule),
		},
	}

	return
}
