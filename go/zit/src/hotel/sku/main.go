package sku

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
)

func init() {
	gob.Register(Transacted{})
	gob.Register(External{})
}

type (
	Ennui interface {
		WriteOneObjekteMetadatei(o *Transacted) (err error)
		ReadOneEnnui(*sha.Sha) (*Transacted, error)
		ReadOneKennung(kennung.Kennung) (*Transacted, error)
		ReadOneKennungSha(kennung.Kennung) (*sha.Sha, error)
	}

	Queryable interface {
		ContainsSku(*Transacted) bool
	}

	Query interface {
		Queryable
		schnittstellen.Stringer
		Each(schnittstellen.FuncIter[Query]) error
	}

	QueryWithSigilAndKennung interface {
		Query
		GetSigil() kennung.Sigil
		ContainsKennung(*kennung.Kennung2) bool
	}

	TransactedAdder interface {
		AddTransacted(*Transacted) error
	}

	QueryGroup interface {
		Query
		Get(gattung.Gattung) (QueryWithSigilAndKennung, bool)
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
