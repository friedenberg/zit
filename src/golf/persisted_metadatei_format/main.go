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
		Format(io.Writer, FormatterContext) (int64, error)
		Parse(io.Reader, ParserContext) (int64, error)
	}
	Getter interface {
		GetPersistentMetadateiFormat() Format
	}
)

func FormatForVersion(v schnittstellen.StoreVersion) Format {
	switch v.Int() {
	case 0:
		return V0{}

	case 1:
		return V1{}

	default:
		return V2{}
	}
}
