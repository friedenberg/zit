package interfaces

import "io"

type (
	SavedBlobFormatter interface {
		FormatSavedBlob(io.Writer, Sha) (int64, error)
	}

	ParseSaver[T any] interface {
		ParseSaveBlob(io.Reader, T) (Sha, int64, error)
	}

	Parser[T any] interface {
		ParseBlob(io.Reader, T) (int64, error)
	}

	Format[T any] interface {
		SavedBlobFormatter
		EncoderTo[T]
		Parser[T]
	}

	TypedBlobStore[T any] interface {
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
