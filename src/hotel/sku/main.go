package sku

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

func init() {
	gob.Register(Transacted{})
	gob.Register(External{})
}

type (
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
		CalculateObjekteSha() (err error)

		EqualsSkuLikePtr(SkuLike) bool

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
