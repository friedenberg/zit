package type_blobs

import (
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
)

const (
	TypeV0     = builtin_types.TypeTypeV0
	TypeV1     = builtin_types.TypeTypeV1
	TypeLatest = builtin_types.TypeTypeLatest
)

func Default() (t TomlV1) {
	t = TomlV1{
		FileExtension: "md",
		Formatters:    make(map[string]script_config.WithOutputFormat),
		VimSyntaxType: "markdown",
	}

	return
}

func DefaultLuaV0() (t V0) {
	t = V0{
		InlineBlob:    true,
		FileExtension: "lua",
		Formatters:    make(map[string]script_config.WithOutputFormat),
		VimSyntaxType: "lua",
		Actions:       make(map[string]script_config.ScriptConfig),
	}

	return
}

func Make() (t *V0) {
	t = &V0{
		Formatters: make(map[string]script_config.WithOutputFormat),
		Actions:    make(map[string]script_config.ScriptConfig),
	}

	return
}
