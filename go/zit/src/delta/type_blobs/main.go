package type_blobs

import (
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
)

const (
	TypeV0     = "toml-type-v0"
	TypeV1     = "toml-type-v1"
	TypeLatest = TypeV1
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
