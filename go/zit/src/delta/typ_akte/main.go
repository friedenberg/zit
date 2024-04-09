package typ_akte

import (
	"code.linenisgreat.com/zit/src/alfa/etikett_rule"
	"code.linenisgreat.com/zit/src/charlie/script_config"
)

func Default() (t V0) {
	t = V0{
		InlineAkte:     true,
		FileExtension:  "md",
		Formatters:     make(map[string]script_config.WithOutputFormat),
		VimSyntaxType:  "markdown",
		Actions:        make(map[string]script_config.ScriptConfig),
		EtikettenRules: make(map[string]etikett_rule.Rule),
	}

	return
}

func DefaultEtikettLuaV0() (t V0) {
	t = V0{
		InlineAkte:     true,
		FileExtension:  "lua",
		Formatters:     make(map[string]script_config.WithOutputFormat),
		VimSyntaxType:  "lua",
		Actions:        make(map[string]script_config.ScriptConfig),
		EtikettenRules: make(map[string]etikett_rule.Rule),
	}

	return
}

func MakeObjekte() (t *V0) {
	t = &V0{
		Formatters:     make(map[string]script_config.WithOutputFormat),
		Actions:        make(map[string]script_config.ScriptConfig),
		EtikettenRules: make(map[string]etikett_rule.Rule),
	}

	return
}
