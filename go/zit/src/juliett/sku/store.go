package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type (
	AbbrStorePresenceGeneric[V any] interface {
		Exists([3]string) error
	}

	AbbrStoreGeneric[V any, VPtr interfaces.Ptr[V]] interface {
		AbbrStorePresenceGeneric[V]
		ExpandStringString(string) (string, error)
		ExpandString(string) (VPtr, error)
		Expand(VPtr) (VPtr, error)
		Abbreviate(ids.Abbreviatable) (string, error)
	}

	AbbrStore interface {
		ZettelId() AbbrStoreGeneric[ids.ZettelId, *ids.ZettelId]
		Shas() AbbrStoreGeneric[sha.Sha, *sha.Sha]

		AddObjectToAbbreviationStore(*Transacted) error
		GetAbbr() ids.Abbr

		errors.Flusher
	}

	ObjectStore interface {
		Commit(ExternalLike, CommitOptions) (err error)
		ReadOneInto(interfaces.ObjectId, *Transacted) (err error)
		ReadPrimitiveQuery(
			qg PrimitiveQueryGroup,
			w interfaces.FuncIter[*Transacted],
		) (err error)
	}

	ExternalObjectId       = ids.ExternalObjectIdLike
	ExternalObjectIdGetter = ids.ExternalObjectIdGetter

	FuncReadOneInto = func(
		k1 interfaces.ObjectId,
		out *Transacted,
	) (err error)

	ExternalStoreUpdateTransacted interface {
		UpdateTransacted(z *Transacted) (err error)
	}

	ExternalStoreReadExternalLikeFromObjectId interface {
		ReadExternalLikeFromObjectId(
			o CommitOptions,
			k1 interfaces.ObjectId,
			t *Transacted,
		) (e ExternalLike, err error)
	}

	ExternalStoreReadAllExternalItems interface {
		ReadAllExternalItems() error
	}

	ExternalStoreForQuery interface {
		GetObjectIdsForString(string) ([]ExternalObjectId, error)
	}

	ExternalStoreForQueryGetter interface {
		GetExternalStoreForQuery(ids.RepoId) (ExternalStoreForQuery, bool)
	}

	BlobStore[T any] interface {
		GetTransactedWithBlob(
			sk TransactedGetter,
		) (common TransactedWithBlob[T], n int64, err error)

		PutTransactedWithBlob(
			TransactedWithBlob[T],
		) error
	}

	BlobCopyResult struct {
		*Transacted    // may be nil
		interfaces.Sha // may not be nil

		// -1: no remote blob store and the blob doesn't exist locally
		// -2: no remote blob store and the blob exists locally
		N int64
	}
)

func MakeBlobCopierDelegate(ui fd.Std) func(BlobCopyResult) error {
	return func(result BlobCopyResult) error {
		return ui.Printf(
			"copied Blob %s (%d bytes)",
			result.Sha,
			result.N,
		)
	}
}
