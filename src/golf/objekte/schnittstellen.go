package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type (
	AkteFormatter interface {
		FormatAkte(io.Writer, schnittstellen.Sha) (int64, error)
	}

	AkteParseSaver[T any] interface {
		ParseSaveAkte(io.Reader, T) (schnittstellen.Sha, int64, error)
	}
)
