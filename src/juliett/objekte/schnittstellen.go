package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/india/matcher"
)

type (
	Konfig interface {
		schnittstellen.Konfig
		matcher.ImplicitEtikettenGetter
		IsInlineTyp(kennung.Typ) bool
	}

	SavedAkteFormatter interface {
		FormatSavedAkte(io.Writer, schnittstellen.ShaLike) (int64, error)
	}

	ParsedAkteFormatter[T any] interface {
		FormatParsedAkte(io.Writer, T) (int64, error)
	}

	AkteParseSaver[T any] interface {
		ParseSaveAkte(io.Reader, T) (schnittstellen.ShaLike, int64, error)
	}

	AkteParser[T any] interface {
		ParseAkte(io.Reader, T) (int64, error)
	}

	AkteFormat[T any, TPtr schnittstellen.Ptr[T]] interface {
		SavedAkteFormatter
		ParsedAkteFormatter[T]
		AkteParser[TPtr]
		// AkteParseSaver[TPtr]
	}
)
