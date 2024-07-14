package sku

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
)

func init() {
	gob.Register(Transacted{})
}

type ObjekteOptions struct {
	ids.RepoId
	objekte_mode.Mode
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
		WriteOneObjekteMetadatei(o *Transacted) (err error)
		ReadOneObjectSha(*sha.Sha) (*Transacted, error)
		ReadOneObjectId(
			interfaces.StringerGenreGetter,
		) (*Transacted, error)
		ReadOneObjectIdSha(
			interfaces.StringerGenreGetter,
		) (*sha.Sha, error)
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
		GetObjekteSha() interfaces.Sha
		GetAkteSha() interfaces.Sha
		GetKey() string

		object_metadata.Getter

		SetBlobSha(interfaces.Sha) error
		SetObjectSha(interfaces.Sha) error
		CalculateObjekteShas() (err error)

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
	}

	CheckedOutLike interface {
		interfaces.Stringer
		TransactedGetter
		ExternalLike
		GetKasten() ids.RepoId
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
		ReadTransactedFromKennung(
			k1 interfaces.StringerGenreGetter,
		) (sk1 *Transacted, err error)
	}
)

func EqualsSkuLike(a, b SkuLike) (ok bool) {
	if !ids.Equals(a.GetObjectId(), b.GetObjectId()) {
		return
	}

	if !a.GetObjekteSha().EqualsSha(b.GetObjekteSha()) {
		return
	}

	return true
}
