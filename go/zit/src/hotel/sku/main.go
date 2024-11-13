package sku

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

func init() {
	gob.Register(Transacted{})
}

// TODO rename and switch to no object.Mode
type CommitOptions struct {
	ids.RepoId
	object_mode.Mode // TODO rename
	ids.Clock
	Proto              *Transacted
	ChangeIsHistorical bool
	DontAddMissingTags bool
	DontAddMissingType bool
	DontValidate       bool
	DontRunHooks       bool
}

type (
	Config interface {
		interfaces.Config
		ids.InlineTypeChecker // TODO move out of konfig entirely
		GetTypeExtension(string) string
	}

	ObjectProbeIndex interface {
		ReadOneObjectId(string, *Transacted) error
	}

	TransactedGetter interface {
		GetSku() *Transacted
	}

	TransactedWithBlob[T any] struct {
		*Transacted
		Blob T
	}

	ExternalLike interface {
		ids.ObjectIdGetter
		interfaces.Stringer
		TransactedGetter
		ExternalLikeGetter
		GetExternalState() external_state.State
		ExternalObjectIdGetter
		GetRepoId() ids.RepoId
	}

	ExternalLikeGetter interface {
		GetSkuExternal() *Transacted
	}

	FSItemReadWriter interface {
		ReadFSItemFromExternal(TransactedGetter) (*FSItem, error)
		WriteFSItemToExternal(*FSItem, TransactedGetter) (err error)
	}

	OneReader interface {
		ReadTransactedFromObjectId(
			k1 interfaces.ObjectId,
		) (sk1 *Transacted, err error)
	}

	BlobSaver interface {
		SaveBlob(ExternalLike) (err error)
	}
)
