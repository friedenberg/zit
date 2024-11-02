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
