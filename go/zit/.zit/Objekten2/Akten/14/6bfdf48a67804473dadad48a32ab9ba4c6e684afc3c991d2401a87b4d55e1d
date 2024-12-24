package blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type (
	ParsedBlobFormatter2[T any] interface {
		FormatParsedBlob(io.Writer, T) (int64, error)
	}

	ParseSaver2[T any] interface {
		ParseSaveBlob(io.Reader, T) (interfaces.Sha, int64, error)
	}

	Parser2[T any] interface {
		ParseBlob(io.Reader, T) (int64, error)
	}

	Format2[T any] interface {
		SavedBlobFormatter
		ParsedBlobFormatter2[T]
		Parser2[T]
	}

	CommonStore[T any] interface {
		ParseTypedBlob(
			tipe interfaces.ObjectId,
			blobSha interfaces.Sha,
		) (common T, n int64, err error)

		PutTypedBlob(
			interfaces.ObjectId,
			T,
		) error
	}

	CommonStore2[T any] interface {
		GetTransactedWithBlob(
			sk sku.TransactedGetter,
		) (common sku.TransactedWithBlob[T], n int64, err error)

		PutTransactedWithBlob(
			sku.TransactedWithBlob[T],
		) error
	}
)
