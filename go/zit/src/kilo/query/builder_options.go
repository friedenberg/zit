package query

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type QueryBuilderModifier interface {
	ModifyBuilder(*Builder)
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

type BuilderOptions interface{}
