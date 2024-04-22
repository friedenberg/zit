package akten

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

type (
	SavedAkteFormatter interface {
		FormatSavedAkte(io.Writer, schnittstellen.ShaLike) (int64, error)
	}

	ParsedAkteFormatter[T any, TPtr schnittstellen.Ptr[T]] interface {
		FormatParsedAkte(io.Writer, TPtr) (int64, error)
	}

	ParseSaver[T any, TPtr schnittstellen.Ptr[T]] interface {
		ParseSaveAkte(io.Reader, TPtr) (schnittstellen.ShaLike, int64, error)
	}

	Parser[T any, TPtr schnittstellen.Ptr[T]] interface {
		ParseAkte(io.Reader, TPtr) (int64, error)
	}

	Format[T any, TPtr schnittstellen.Ptr[T]] interface {
		SavedAkteFormatter
		ParsedAkteFormatter[T, TPtr]
		Parser[T, TPtr]
	}

	Konfig interface {
		schnittstellen.Konfig
		IsInlineTyp(kennung.Typ) bool
		GetApproximatedTyp(kennung.Kennung) ApproximatedTyp
	}
)
