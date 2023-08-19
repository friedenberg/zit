package typ

import (
	"github.com/friedenberg/zit/src/alfa/etikett_rule"
	"github.com/friedenberg/zit/src/charlie/script_config"
	"github.com/friedenberg/zit/src/echo/kennung"
)

func Default() (t Akte, k kennung.Typ) {
	k = kennung.MustTyp("md")

	t = Akte{
		InlineAkte:     true,
		FileExtension:  "md",
		Formatters:     make(map[string]script_config.WithOutputFormat),
		VimSyntaxType:  "markdown",
		Actions:        make(map[string]script_config.ScriptConfig),
		EtikettenRules: make(map[string]etikett_rule.Rule),
	}

	return
}

func MakeObjekte() (t *Akte) {
	t = &Akte{
		Formatters:     make(map[string]script_config.WithOutputFormat),
		Actions:        make(map[string]script_config.ScriptConfig),
		EtikettenRules: make(map[string]etikett_rule.Rule),
	}

	return
}
