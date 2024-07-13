package blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type (
	SavedBlobFormatter interface {
		FormatSavedBlob(io.Writer, interfaces.ShaLike) (int64, error)
	}

	ParsedBlobFormatter[T any, TPtr interfaces.Ptr[T]] interface {
		FormatParsedBlob(io.Writer, TPtr) (int64, error)
	}

	ParseSaver[T any, TPtr interfaces.Ptr[T]] interface {
		ParseSaveBlob(io.Reader, TPtr) (interfaces.ShaLike, int64, error)
	}

	Parser[T any, TPtr interfaces.Ptr[T]] interface {
		ParseBlob(io.Reader, TPtr) (int64, error)
	}

	Format[T any, TPtr interfaces.Ptr[T]] interface {
		SavedBlobFormatter
		ParsedBlobFormatter[T, TPtr]
		Parser[T, TPtr]
	}

	Config interface {
		interfaces.Config
		IsInlineTyp(ids.Type) bool
	}
)
