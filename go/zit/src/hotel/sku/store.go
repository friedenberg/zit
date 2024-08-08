package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type (
	ExternalObjectId interface {
		ids.ObjectIdLike
		GetExternalObjectId() *ids.ObjectId
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

	ExternalStoreForQuery interface {
		GetObjectIdsForString(string) ([]ExternalObjectId, error)
	}

	ExternalStoreForQueryGetter interface {
		GetExternalStoreForQuery(ids.RepoId) (ExternalStoreForQuery, bool)
	}
)
