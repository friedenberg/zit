package sku

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/golf/objekte_format"
)

func init() {
	gob.Register(Transacted{})
}

type ObjekteOptions struct {
	kennung.RepoId
	objekte_mode.Mode
	kennung.Clock
	Proto *Transacted
}

type (
	Konfig interface {
		interfaces.Config
		kennung.InlineTypChecker // TODO move out of konfig entirely
		GetTypExtension(string) string
	}

	Ennui interface {
		WriteOneObjekteMetadatei(o *Transacted) (err error)
		ReadOneEnnui(*sha.Sha) (*Transacted, error)
		ReadOneKennung(
			interfaces.StringerGenreGetter,
		) (*Transacted, error)
		ReadOneKennungSha(
			interfaces.StringerGenreGetter,
		) (*sha.Sha, error)
	}

	TransactedAdder interface {
		AddTransacted(*Transacted) error
	}

	SkuLike interface {
		interfaces.ValueLike
		interfaces.Stringer
		metadatei.Getter

		GetKopf() kennung.Tai
		GetTai() kennung.Tai
		GetTyp() kennung.Type
		GetKennung() kennung.IdLike
		GetObjekteSha() interfaces.ShaLike
		GetAkteSha() interfaces.ShaLike
		GetKey() string

		metadatei.Getter

		SetAkteSha(interfaces.ShaLike) error
		SetObjekteSha(interfaces.ShaLike) error
		CalculateObjekteShas() (err error)

		SetTai(kennung.Tai)
		SetKennungLike(kennung.IdLike) error
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
		GetKasten() kennung.RepoId
		GetSkuCheckedOutLike() CheckedOutLike
		GetState() checked_out_state.State
		SetState(checked_out_state.State) error
		GetError() error
		Clone() CheckedOutLike
	}

	ManyPrinter interface {
		PrintMany(...objekte_format.FormatterContext) (int64, error)
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
	if !kennung.Equals(a.GetKennung(), b.GetKennung()) {
		return
	}

	if !a.GetObjekteSha().EqualsSha(b.GetObjekteSha()) {
		return
	}

	return true
}
