package objekte_format

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
)

type (
	FormatterContext interface {
		metadatei.PersistentFormatterContext
		GetKennung() kennung.Id
	}

	ParserContext interface {
		metadatei.PersistentParserContext
		SetKennungLike(kennung.Id) error
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

func FormatForVersion(sv interfaces.StoreVersion) Format {
	v := sv.GetInt()

	switch v {
	case 3, 4:
		return v4{}

	default:
		return v5{}
	}
}

func FormatForVersions(write, read interfaces.StoreVersion) Format {
	return MakeBespoke(
		FormatForVersion(write),
		FormatForVersion(read),
	)
}

type nopFormatterContext struct {
	metadatei.PersistentFormatterContext
}

func (nopFormatterContext) GetKennung() kennung.Id {
	return nil
}
