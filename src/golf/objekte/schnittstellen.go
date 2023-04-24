package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type (
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

	VerzeichnissePtr[T any, T1 Objekte[T1]] interface {
		schnittstellen.Resetable[T]
		ResetWithObjekteMetadateiGetter(T1, metadatei.Getter)
	}
)
