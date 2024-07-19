package sku

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
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

	TransactedAdder interface {
		AddTransacted(*Transacted) error
	}

	SkuLike interface {
		interfaces.ValueLike
		interfaces.Stringer
		object_metadata.Getter

		GetTai() ids.Tai
		GetType() ids.Type
		GetObjectId() *ids.ObjectId
		GetObjectSha() interfaces.Sha
		GetBlobSha() interfaces.Sha
		GetKey() string

		object_metadata.Getter

		SetBlobSha(interfaces.Sha) error
		SetObjectSha(interfaces.Sha) error
		CalculateObjectShas() (err error)

		SetTai(ids.Tai)
		object_inventory_format.ParserContext
		object_inventory_format.FormatterContext
		SetFromSkuLike(SkuLike) error

		GetSkuLike() SkuLike
	}

	TransactedGetter interface {
		GetSku() *Transacted
	}

	ExternalLikeGetter interface {
		GetSkuExternalLike() ExternalLike
	}

	ExternalLike interface {
		interfaces.Stringer
		TransactedGetter
		ExternalLikeGetter
		Clone() ExternalLike
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

func EqualsSkuLike(a, b SkuLike) (ok bool) {
	if !ids.Equals(a.GetObjectId(), b.GetObjectId()) {
		return
	}

	if !a.GetObjectSha().EqualsSha(b.GetObjectSha()) {
		return
	}

	return true
}
