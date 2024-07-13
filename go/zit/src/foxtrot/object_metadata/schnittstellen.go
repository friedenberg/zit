package object_metadata

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
)

type (
	Getter interface {
		GetMetadata() *Metadata
	}

	Setter interface {
		SetMetadata(*Metadata)
	}

	MetadataLike interface {
		Getter
	}

	BlobPathGetter interface {
		GetBlobPath() string
	}

	BlobPathSetter interface {
		SetBlobFD(*fd.FD) error
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
		SetBlobSha(interfaces.Sha) error
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
