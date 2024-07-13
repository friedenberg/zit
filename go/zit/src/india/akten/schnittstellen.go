package akten

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type (
	SavedAkteFormatter interface {
		FormatSavedAkte(io.Writer, interfaces.ShaLike) (int64, error)
	}

	ParsedAkteFormatter[T any, TPtr interfaces.Ptr[T]] interface {
		FormatParsedAkte(io.Writer, TPtr) (int64, error)
	}

	ParseSaver[T any, TPtr interfaces.Ptr[T]] interface {
		ParseSaveAkte(io.Reader, TPtr) (interfaces.ShaLike, int64, error)
	}

	Parser[T any, TPtr interfaces.Ptr[T]] interface {
		ParseAkte(io.Reader, TPtr) (int64, error)
	}

	Format[T any, TPtr interfaces.Ptr[T]] interface {
		SavedAkteFormatter
		ParsedAkteFormatter[T, TPtr]
		Parser[T, TPtr]
	}

	Konfig interface {
		interfaces.Konfig
		IsInlineTyp(kennung.Typ) bool
	}
)
