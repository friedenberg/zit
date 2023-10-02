package objekte_format

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type (
	FormatterContext interface {
		metadatei.PersistentFormatterContext
		GetKennungLike() kennung.Kennung
	}

	ParserContext interface {
		metadatei.PersistentParserContext
		SetKennungLike(kennung.Kennung) error
	}

	Formatter interface {
		FormatPersistentMetadatei(
			io.Writer,
			FormatterContext,
			Options,
		) (int64, error)
	}

	Parser interface {
		ParsePersistentMetadatei(
			io.Reader,
			ParserContext,
			Options,
		) (int64, error)
	}

	Format interface {
		Formatter
		Parser
	}

	Getter interface {
		GetPersistentMetadateiFormat() Format
	}
)

func Default() Format {
	return v4{}
}

func FormatForVersion(v schnittstellen.StoreVersion) Format {
	switch v.GetInt() {
	case 0:
		return v0{}

	case 1:
		return v1{}

	case 2:
		return v2{}

	case 3:
		fallthrough
		// return v3{}

	default:
		return v4{}
	}
}

func FormatForVersions(write, read schnittstellen.StoreVersion) Format {
	return MakeBespoke(
		FormatForVersion(write),
		FormatForVersion(read),
	)
}
