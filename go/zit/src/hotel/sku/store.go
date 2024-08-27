package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type (
	ExternalObjectId interface {
		interfaces.GenreGetter
		interfaces.Stringer
		ExternalObjectIdGetter
	}

	ExternalObjectIdGetter interface {
		GetExternalObjectId() ExternalObjectId
	}

	FuncRealize     = func(*Transacted, *Transacted, CommitOptions) error
	FuncCommit      = func(*Transacted, CommitOptions) error
	FuncReadSha     = func(*sha.Sha) (*Transacted, error)
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

	ExternalStoreApplyDotOperator interface {
		ApplyDotOperator() error
	}

	ExternalStoreForQuery interface {
		GetObjectIdsForString(string) ([]ExternalObjectId, error)
	}

	ExternalStoreForQueryGetter interface {
		GetExternalStoreForQuery(ids.RepoId) (ExternalStoreForQuery, bool)
	}
)
