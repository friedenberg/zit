package sku

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type (
	Metadatei = metadatei.Metadatei
)

type Kennung = kennung.Kennung

type IdLikeGetter interface {
	GetId() schnittstellen.ValueLike
}

type DataIdentity interface {
	schnittstellen.ValueLike
	GetTai() kennung.Tai
	GetId() Kennung
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

	GetTransactionIndex() values.Int
	// Less(SkuLike) bool
}

type SkuLikePtr interface {
	SkuLike
	SetTimeAndFields(kennung.Tai, ...string) error
	SetFromSku(Sku) error
	SetTransactionIndex(int)
}
