package sku

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
)

func init() {
	gob.Register(Transacted{})
}

type ObjekteOptions struct {
	objekte_mode.Mode
	kennung.Clock
}

type (
	Konfig interface {
		schnittstellen.Konfig
		kennung.InlineTypChecker // TODO move out of konfig entirely
		GetTypExtension(string) string
	}

	FuncRealize = func(*Transacted, *Transacted, ObjekteOptions) error
	FuncCommit  = func(*Transacted, ObjekteOptions) error
	FuncReadSha = func(*sha.Sha) (*Transacted, error)

	StoreFuncs struct {
		FuncRealize
		FuncCommit
		FuncReadSha
	}

	Ennui interface {
		WriteOneObjekteMetadatei(o *Transacted) (err error)
		ReadOneEnnui(*sha.Sha) (*Transacted, error)
		ReadOneKennung(
			schnittstellen.StringerGattungGetter,
		) (*Transacted, error)
		ReadOneKennungSha(
			schnittstellen.StringerGattungGetter,
		) (*sha.Sha, error)
	}

	TransactedAdder interface {
		AddTransacted(*Transacted) error
	}

	SkuLike interface {
		schnittstellen.ValueLike
		schnittstellen.Stringer
		metadatei.Getter

		GetKopf() kennung.Tai
		GetTai() kennung.Tai
		GetTyp() kennung.Typ
		GetKennung() kennung.Kennung
		GetObjekteSha() schnittstellen.ShaLike
		GetAkteSha() schnittstellen.ShaLike
		GetKey() string

		metadatei.Getter

		SetAkteSha(schnittstellen.ShaLike) error
		SetObjekteSha(schnittstellen.ShaLike) error
		CalculateObjekteShas() (err error)

		SetTai(kennung.Tai)
		SetKennungLike(kennung.Kennung) error
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
		schnittstellen.Stringer
		TransactedGetter
		ExternalLikeGetter
	}

	CheckedOutLike interface {
		schnittstellen.Stringer
		TransactedGetter
		ExternalLike
		GetSkuCheckedOutLike() CheckedOutLike
		GetState() checked_out_state.State
		SetState(checked_out_state.State) error
		GetError() error
	}

	OneReader interface {
		ReadOne(
			k1 schnittstellen.StringerGattungGetter,
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
