package sku

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type (
	Metadatei = metadatei.Metadatei
	Kennung   = kennung.Kennung

	IdLikeGetter interface {
		GetId() schnittstellen.ValueLike
	}
)

type Getter interface {
	GetSkuLike() SkuLike
}

type WithKennungInterface interface {
	SkuLike
}

type SkuLike interface {
	schnittstellen.ValueLike
	schnittstellen.GattungGetter
	metadatei.Getter

	GetTai() kennung.Tai
	GetId() Kennung
	GetKennungLike() kennung.Kennung
	GetObjekteSha() schnittstellen.ShaLike
	GetAkteSha() schnittstellen.ShaLike
	GetKey() string

	EqualsSkuLike(SkuLike) bool
	ImmutableClone() SkuLike
	MutableClone() SkuLikePtr
}

type SkuLikePtr interface {
	SkuLike

	metadatei.GetterPtr
	metadatei.Setter

	SetKennungLike(kennung.Kennung) error
	GetKennungLikePtr() kennung.KennungPtr
	SetFromSkuLike(SkuLike) error
	Reset()
}
