package sku

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
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
		CloneExternalLike() ExternalLike
		GetExternalState() external_state.State
		ExternalObjectIdGetter
		GetRepoId() ids.RepoId
	}

	ExternalLikeGetter interface {
		GetSkuExternalLike() ExternalLike
	}

	FSItemReadWriter interface {
		ReadFSItemFromExternal(el ExternalLike) (i *FSItem, err error)
		WriteFSItemToExternal(i *FSItem, el ExternalLike) (err error)
	}

	CheckedOutLike interface {
		interfaces.Stringer
		TransactedGetter
		ExternalLikeGetter
		GetSkuCheckedOutLike() CheckedOutLike
		GetState() checked_out_state.State
		SetState(checked_out_state.State) error
		GetError() error
		CloneCheckedOutLike() CheckedOutLike
	}

	ManyPrinter interface {
		PrintMany(...object_inventory_format.FormatterContext) (int64, error)
	}

	Scanner interface {
		Scan() bool
		GetTransacted() *Transacted
		Error() error
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
