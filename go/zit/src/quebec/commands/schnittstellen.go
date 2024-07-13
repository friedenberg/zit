package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

type QueryBuilderModifier interface {
	ModifyBuilder(*query.Builder)
}

type DefaultSigilGetter interface {
	DefaultSigil() kennung.Sigil
}

type DefaultGattungGetter interface {
	DefaultGattungen() kennung.Genre
}

type CompletionGattungGetter interface {
	CompletionGattung() kennung.Genre
}
