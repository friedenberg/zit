package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type (
	ExternalObjectId       = ids.ExternalObjectId
	ExternalObjectIdGetter = ids.ExternalObjectIdGetter

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

	ExternalStoreReadAllExternalItems interface {
		ReadAllExternalItems() error
	}

	ExternalStoreForQuery interface {
		GetObjectIdsForString(string) ([]ExternalObjectId, error)
	}

	ExternalStoreForQueryGetter interface {
		GetExternalStoreForQuery(ids.RepoId) (ExternalStoreForQuery, bool)
	}

	ExternalLikePoolGetter interface {
		GetExternalLikePool() interfaces.PoolValue[ExternalLike]
	}

	ExternalLikeResetter3Getter interface {
		GetExternalLikeResetter3() interfaces.Resetter3[ExternalLike]
	}
)
