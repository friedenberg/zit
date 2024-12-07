package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

type QueryBuilderModifier interface {
	ModifyBuilder(*query.Builder)
}

type DefaultSigilGetter interface {
	DefaultSigil() ids.Sigil
}

type DefaultGenresGetter interface {
	DefaultGenres() ids.Genre
}

type CompletionGenresGetter interface {
	CompletionGenres() ids.Genre
}
