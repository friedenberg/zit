package metadatei

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type (
	Getter interface {
		GetMetadatei() Metadatei
	}

	Setter interface {
		SetMetadatei(Metadatei)
	}

	AktePathGetter interface {
		GetAktePath() string
	}

	AktePathSetter interface {
		SetAkteFD(kennung.FD) error
	}

	PersistentFormatterContext interface {
		Getter
		GetAkteSha() schnittstellen.Sha
	}

	PersistentParserContext interface {
		Getter
		Setter
	}

	TextFormatterContext interface {
		PersistentFormatterContext
		// GetAktePath() string
	}

	TextParserContext interface {
		PersistentParserContext
	}

	TextFormatOutput struct {
		io.Writer
		string
	}

	TextFormatter interface {
		Format(io.Writer, TextFormatterContext) (int64, error)
	}

	TextParser interface {
		Parse(io.Reader, TextParserContext) (int64, error)
	}
)
