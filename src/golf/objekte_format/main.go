package objekte_format

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type (
	FormatterContext interface {
		metadatei.PersistentFormatterContext
		GetKennung() kennung.Kennung
	}

	FormatterContextIncludeTai interface {
		FormatterContext
		IncludeTai() bool
	}

	ParserContext interface {
		metadatei.PersistentParserContext
		SetKennung(kennung.Kennung) error
	}

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

func BestandsaufnahmeFormatIncludeTai() Format {
	return v3{includeTai: true}
}

func BestandsaufnahmeFormatExcludeTai() Format {
	return v3{}
}

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
