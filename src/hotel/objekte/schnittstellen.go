package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type (
	Konfig interface {
		schnittstellen.Konfig
		kennung.ImplicitEtikettenGetter
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

	AkteFormat[T any, TPtr schnittstellen.Ptr[T]] interface {
		SavedAkteFormatter
		ParsedAkteFormatter[T]
		AkteParseSaver[TPtr]
	}

	VerzeichnissePtr[T any, T1 Akte[T1]] interface {
		schnittstellen.Resetable[T]
		ResetWithObjekteMetadateiGetter(T1, metadatei.Getter)
	}
)