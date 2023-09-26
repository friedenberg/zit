package sku

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

func init() {
	gob.Register(Transacted{})
	gob.Register(External{})
}

type (
	skuLike interface {
		schnittstellen.ValueLike
		schnittstellen.GattungGetter
		metadatei.Getter

		GetKopf() kennung.Tai
		GetTai() kennung.Tai
		GetTyp() kennung.Typ
		GetKennungLike() kennung.Kennung
		GetObjekteSha() schnittstellen.ShaLike
		GetAkteSha() schnittstellen.ShaLike
		GetKey() string
	}

	SkuLikePtr interface {
		skuLike

		metadatei.GetterPtr
		metadatei.Setter

		SetAkteSha(schnittstellen.ShaLike)
		SetObjekteSha(schnittstellen.ShaLike)

		EqualsSkuLikePtr(SkuLikePtr) bool

		SetTai(kennung.Tai)
		SetKennungLike(kennung.Kennung) error
		GetKennungLikePtr() kennung.KennungPtr
		SetFromSkuLike(SkuLikePtr) error
		Reset()
	}

	SkuLikeExternalPtr interface {
		SkuLikePtr

		GetFDs() ExternalFDs
		GetFDsPtr() *ExternalFDs
		GetAkteFD() kennung.FD
		GetAktePath() string

		GetObjekteFD() kennung.FD

		ResetWithExternalMaybe(b ExternalMaybe) (err error)
	}
)
