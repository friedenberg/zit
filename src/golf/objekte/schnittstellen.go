package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type (
	// AkteGetter[A any] interface {
	// 	GetAkte() A
	// }

	// AktePtrGetter[A any, APtr schnittstellen.Ptr[A]] interface {
	// 	GetAktePtr() APtr
	// }

	AkteParser[T any] interface {
		ParseAkte(io.Reader, T) (int64, error)
	}

	SavedAkteFormatter interface {
		FormatSavedAkte(io.Writer, schnittstellen.Sha) (int64, error)
	}

	ParsedAkteFormatter[T any] interface {
		FormatParsedAkte(io.Writer, T) (int64, error)
	}

	AkteParseSaver[T any] interface {
		ParseSaveAkte(io.Reader, T) (schnittstellen.Sha, int64, error)
	}

	AkteFormat[T any, TPtr schnittstellen.Ptr[T]] interface {
		SavedAkteFormatter
		ParsedAkteFormatter[T]
		AkteParseSaver[TPtr]
	}
)
