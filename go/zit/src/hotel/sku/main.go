package sku

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
)

func init() {
	gob.Register(Transacted{})
}

type CommitOptions struct {
	ids.RepoId
	objekte_mode.Mode // TODO rename
	ids.Clock
	Proto *Transacted
}

type (
	Config interface {
		interfaces.Config
		ids.InlineTypeChecker // TODO move out of konfig entirely
		GetTypeExtension(string) string
	}

	ObjectProbeIndex interface {
		ReadOneObjectId(
			interfaces.ObjectId,
		) (*Transacted, error)
	}

	TransactedGetter interface {
		GetSku() *Transacted
	}

	ExternalLike interface {
		ids.ObjectIdGetter
		interfaces.Stringer
		TransactedGetter
		ExternalLikeGetter
		Clone() ExternalLike
	}

	ExternalLikeGetter interface {
		GetSkuExternalLike() ExternalLike
	}

	CheckedOutLike interface {
		interfaces.Stringer
		TransactedGetter
		ExternalLikeGetter
		GetRepoId() ids.RepoId
		GetSkuCheckedOutLike() CheckedOutLike
		GetState() checked_out_state.State
		SetState(checked_out_state.State) error
		GetError() error
		Clone() CheckedOutLike
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
)
