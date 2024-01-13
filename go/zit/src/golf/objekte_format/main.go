package objekte_format

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type (
	FormatterContext interface {
		metadatei.PersistentFormatterContext
		GetKennung() kennung.Kennung
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
			*catgut.RingBuffer,
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

func FormatForVersion(sv schnittstellen.StoreVersion) Format {
	v := sv.GetInt()

	switch v {
	case 3, 4:
		return v4{}

	default:
		return v5{}
	}
}

func FormatForVersions(write, read schnittstellen.StoreVersion) Format {
	return MakeBespoke(
		FormatForVersion(write),
		FormatForVersion(read),
	)
}

type nopFormatterContext struct {
	metadatei.PersistentFormatterContext
}

func (nopFormatterContext) GetKennung() kennung.Kennung {
	return nil
}
