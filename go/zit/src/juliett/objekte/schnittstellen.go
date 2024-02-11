package objekte

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/india/matcher"
)

type (
	Konfig interface {
		schnittstellen.Konfig
		matcher.ImplicitEtikettenGetter
		IsInlineTyp(kennung.Typ) bool
		GetApproximatedTyp(kennung.Kennung) ApproximatedTyp
	}

	SavedAkteFormatter interface {
		FormatSavedAkte(io.Writer, schnittstellen.ShaLike) (int64, error)
	}

	ParsedAkteFormatter[T any, TPtr schnittstellen.Ptr[T]] interface {
		FormatParsedAkte(io.Writer, TPtr) (int64, error)
	}

	AkteParseSaver[T any, TPtr schnittstellen.Ptr[T]] interface {
		ParseSaveAkte(io.Reader, TPtr) (schnittstellen.ShaLike, int64, error)
	}

	AkteParser[T any, TPtr schnittstellen.Ptr[T]] interface {
		ParseAkte(io.Reader, TPtr) (int64, error)
	}

	AkteFormat[T any, TPtr schnittstellen.Ptr[T]] interface {
		SavedAkteFormatter
		ParsedAkteFormatter[T, TPtr]
		AkteParser[T, TPtr]
	}
)
