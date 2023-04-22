package persisted_metadatei_format

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type (
	FormatterContext = metadatei.PersistentFormatterContext
	ParserContext    = metadatei.PersistentParserContext
	Format           interface {
		FormatPersistentMetadatei(io.Writer, FormatterContext) (int64, error)
		ParsePersistentMetadatei(io.Reader, ParserContext) (int64, error)
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
