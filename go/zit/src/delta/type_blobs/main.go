package type_blobs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/etikett_rule"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
)

func Default() (t V0) {
	t = V0{
		InlineBlob:     true,
		FileExtension:  "md",
		Formatters:     make(map[string]script_config.WithOutputFormat),
		VimSyntaxType:  "markdown",
		Actions:        make(map[string]script_config.ScriptConfig),
		EtikettenRules: make(map[string]etikett_rule.Rule),
	}

	return
}

func DefaultLuaV0() (t V0) {
	t = V0{
		InlineBlob:     true,
		FileExtension:  "lua",
		Formatters:     make(map[string]script_config.WithOutputFormat),
		VimSyntaxType:  "lua",
		Actions:        make(map[string]script_config.ScriptConfig),
		EtikettenRules: make(map[string]etikett_rule.Rule),
	}

	return
}

func Make() (t *V0) {
	t = &V0{
		Formatters:     make(map[string]script_config.WithOutputFormat),
		Actions:        make(map[string]script_config.ScriptConfig),
		EtikettenRules: make(map[string]etikett_rule.Rule),
	}

	return
}
