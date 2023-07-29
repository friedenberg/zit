package sku

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
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
		GetKennungLike() kennung.Kennung
		GetObjekteSha() schnittstellen.ShaLike
		GetAkteSha() schnittstellen.ShaLike
		GetKey() string

		EqualsSkuLike(SkuLike) bool
		ImmutableClone() SkuLike
		MutableClone() SkuLikePtr
	}

	SkuLikePtr interface {
		SkuLike

		metadatei.GetterPtr
		metadatei.Setter

		SetObjekteSha(schnittstellen.ShaLike)

		SetKennungLike(kennung.Kennung) error
		GetKennungLikePtr() kennung.KennungPtr
		SetFromSkuLike(SkuLike) error
		Reset()
	}
)
