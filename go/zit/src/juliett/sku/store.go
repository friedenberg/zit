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

	ExternalStoreReadExternalLikeFromObjectIdLike interface {
		ReadExternalLikeFromObjectIdLike(
			o CommitOptions,
			oid interfaces.Stringer,
			t *Transacted,
		) (e ExternalLike, err error)
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

	ImporterOptions struct {
		ExcludeObjects      bool
		RemoteBlobStore     interfaces.BlobStore
		PrintCopies         bool
		AllowMergeConflicts bool
		BlobCopierDelegate  interfaces.FuncIter[BlobCopyResult]
		ParentNegotiator    ParentNegotiator
		CheckedOutPrinter   interfaces.FuncIter[*CheckedOut]
	}

	Importer interface {
		GetCheckedOutPrinter() interfaces.FuncIter[*CheckedOut]

		SetCheckedOutPrinter(
			p interfaces.FuncIter[*CheckedOut],
		)

		ImportBlobIfNecessary(
			sk *Transacted,
		) (err error)

		Import(
			external *Transacted,
		) (co *CheckedOut, err error)
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
