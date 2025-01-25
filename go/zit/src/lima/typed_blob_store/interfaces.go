package typed_blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type (
	Format[T any, TPtr interfaces.Ptr[T]] interface {
		interfaces.SavedBlobFormatter
		interfaces.Coder[TPtr]
	}
)

type TypedStore[
	A interfaces.Blob[A],
	APtr interfaces.BlobPtr[A],
] interface {
	SaveBlobText(APtr) (interfaces.Sha, int64, error)
	Format[A, APtr]
	interfaces.BlobPool[APtr]
}
