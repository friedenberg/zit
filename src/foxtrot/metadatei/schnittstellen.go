package metadatei

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/fd"
)

type (
	Getter interface {
		GetMetadatei() Metadatei
		GetMetadateiPtr() *Metadatei
	}

	GetterPtr interface {
	}

	Setter interface {
		SetMetadatei(Metadatei)
	}

	MetadateiLike interface {
		Getter
		Setter
	}

	AktePathGetter interface {
		GetAktePath() string
	}

	AktePathSetter interface {
		SetAkteFD(fd.FD) error
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
		SetAkteSha(schnittstellen.ShaLike)
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
