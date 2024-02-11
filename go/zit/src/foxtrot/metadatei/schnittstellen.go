package metadatei

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/echo/fd"
)

type (
	Getter interface {
		GetMetadatei() *Metadatei
	}

	Setter interface {
		SetMetadatei(*Metadatei)
	}

	MetadateiLike interface {
		Getter
	}

	AktePathGetter interface {
		GetAktePath() string
	}

	AktePathSetter interface {
		SetAkteFD(*fd.FD) error
	}

	PersistentFormatterContext interface {
		Getter
	}

	PersistentParserContext interface {
		Getter
	}

	TextFormatterContext interface {
		PersistentFormatterContext
		// GetAktePath() string
	}

	TextParserContext interface {
		PersistentParserContext
		SetAkteSha(schnittstellen.ShaLike) error
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
