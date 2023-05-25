package objekte_format

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type (
	FormatterContext = metadatei.PersistentFormatterContext

	FormatterContextIncludeTai interface {
		FormatterContext
		IncludeTai() bool
	}

	ParserContext = metadatei.PersistentParserContext

	Formatter interface {
		FormatPersistentMetadatei(io.Writer, FormatterContext) (int64, error)
	}

	Parser interface {
		ParsePersistentMetadatei(io.Reader, ParserContext) (int64, error)
	}

	Format interface {
		Formatter
		Parser
	}

	Getter interface {
		GetPersistentMetadateiFormat() Format
	}
)

func FormatForVersion(v schnittstellen.StoreVersion) Format {
	switch v.Int() {
	case 0:
		return v0{}

	case 1:
		return v1{}

	default:
		return v2{}
	}
}

func FormatForVersions(write, read schnittstellen.StoreVersion) Format {
	return MakeBespoke(
		FormatForVersion(write),
		FormatForVersion(read),
	)
}
