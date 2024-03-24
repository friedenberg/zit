package sku

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
)

func init() {
	gob.Register(Transacted{})
	gob.Register(External{})
}

type (
	QueryBase interface {
		ContainsMatchable(*Transacted) bool
		schnittstellen.Stringer
		Each(schnittstellen.FuncIter[QueryBase]) error
	}

	Query interface {
		QueryBase
		GetSigil() kennung.Sigil
		GetKennungen() map[string]*kennung.Kennung2
	}

	MatchableAdder interface {
		AddMatchable(*Transacted) error
	}

	MatcherGroup interface {
		QueryBase
		Get(gattung.Gattung) (Query, bool)
	}

	SkuLike interface {
		schnittstellen.ValueLike
		schnittstellen.GattungGetter
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

	SkuExternalLike interface {
		SkuLike

		GetExternalSkuLike() SkuExternalLike

		GetFDs() *ExternalFDs
		GetAkteFD() *fd.FD
		GetObjekteFD() *fd.FD
		GetAktePath() string

		ResetWithExternalMaybe(b *ExternalMaybe) (err error)
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
