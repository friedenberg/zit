package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

type QueryBuilderModifier interface {
	ModifyBuilder(*query.Builder)
}

type DefaultSigilGetter interface {
	DefaultSigil() ids.Sigil
}

type DefaultGattungGetter interface {
	DefaultGattungen() ids.Genre
}

type CompletionGattungGetter interface {
	CompletionGattung() ids.Genre
}
