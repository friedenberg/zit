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
