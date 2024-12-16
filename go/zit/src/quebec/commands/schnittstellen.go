package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type CompletionGenresGetter interface {
	CompletionGenres() ids.Genre
}
