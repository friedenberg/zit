package blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type (
	ParsedBlobFormatter[T any, TPtr interfaces.Ptr[T]] interface {
		FormatParsedBlob(io.Writer, TPtr) (int64, error)
	}

	ParseSaver[T any, TPtr interfaces.Ptr[T]] interface {
		ParseSaveBlob(io.Reader, TPtr) (interfaces.Sha, int64, error)
	}

	Parser[T any, TPtr interfaces.Ptr[T]] interface {
		ParseBlob(io.Reader, TPtr) (int64, error)
	}

	Format[T any, TPtr interfaces.Ptr[T]] interface {
		interfaces.SavedBlobFormatter
		ParsedBlobFormatter[T, TPtr]
		Parser[T, TPtr]
	}
)

type Store[
	A interfaces.Blob[A],
	APtr interfaces.BlobPtr[A],
] interface {
	SaveBlobText(APtr) (interfaces.Sha, int64, error)
	Format[A, APtr]
	interfaces.BlobGetterPutter[APtr]
}
