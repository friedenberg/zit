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

type DataIdentity interface {
	schnittstellen.ValueLike
	GetTai() kennung.Tai
	GetId() Kennung
	GetKennungLike() kennung.Kennung
	schnittstellen.GattungGetter
	GetObjekteSha() schnittstellen.ShaLike
	GetAkteSha() schnittstellen.ShaLike
	metadatei.Getter
}

type DataIdentityGetter interface {
	GetDataIdentity() DataIdentity
}

type SkuLike interface {
	DataIdentity
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
	// SetTimeAndFields(kennung.Tai, ...string) error
	SetFromSkuLike(SkuLike) error
	// SetTransactionIndex(int)
	Reset()
}
