package type_blobs

import (
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
)

const (
	TypeV0     = builtin_types.TypeTypeTomlV0
	TypeV1     = builtin_types.TypeTypeTomlV1
	TypeLatest = builtin_types.TypeTypeLatestDefault
)

func Default() (t TomlV1) {
	t = TomlV1{
		FileExtension: "md",
		Formatters:    make(map[string]script_config.WithOutputFormat),
		VimSyntaxType: "markdown",
	}

	return
}
