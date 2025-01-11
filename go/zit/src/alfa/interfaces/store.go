package interfaces

import "io"

type (
	SavedBlobFormatter interface {
		FormatSavedBlob(io.Writer, Sha) (int64, error)
	}

	ParsedBlobFormatter[T any] interface {
		FormatParsedBlob(io.Writer, T) (int64, error)
	}

	ParseSaver[T any] interface {
		ParseSaveBlob(io.Reader, T) (Sha, int64, error)
	}

	Parser[T any] interface {
		ParseBlob(io.Reader, T) (int64, error)
	}

	Format[T any] interface {
		SavedBlobFormatter
		ParsedBlobFormatter[T]
		Parser[T]
	}

	CommonStore[T any] interface {
		ParseTypedBlob(
			tipe ObjectId,
			blobSha Sha,
		) (common T, n int64, err error)

		PutTypedBlob(
			ObjectId,
			T,
		) error
	}
)
