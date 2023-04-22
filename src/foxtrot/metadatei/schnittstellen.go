package metadatei

import (
	"io"

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
		FormatMetadatei(io.Writer, TextFormatterContext) (int64, error)
	}

	TextParser interface {
		ParseMetadatei(io.Reader, TextParserContext) (int64, error)
	}
)
