package sku

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type (
	IdLikeGetter interface {
		GetId() schnittstellen.ValueLike
	}

	Getter interface {
		GetSkuLike() SkuLike
		// GetSkuLikePtr() SkuLikePtr
	}

	SkuLike interface {
		schnittstellen.ValueLike
		schnittstellen.GattungGetter
		metadatei.Getter

		GetTai() kennung.Tai
		GetTyp() kennung.Typ
		GetKennungLike() kennung.Kennung
		GetObjekteSha() schnittstellen.ShaLike
		GetAkteSha() schnittstellen.ShaLike
		GetKey() string

		EqualsSkuLike(SkuLike) bool
		ImmutableClone() SkuLike
		MutableClone() SkuLikePtr

		GetSkuLike() SkuLike
	}

	SkuLikePtr interface {
		SkuLike

		metadatei.GetterPtr
		metadatei.Setter

		SetObjekteSha(schnittstellen.ShaLike)

		SetTai(kennung.Tai)
		SetKennungLike(kennung.Kennung) error
		GetKennungLikePtr() kennung.KennungPtr
		SetFromSkuLike(SkuLike) error
		Reset()

		GetSkuLikePtr() SkuLikePtr
	}
)