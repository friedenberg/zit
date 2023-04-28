package sku

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/echo/ts"
)

type IdLike = schnittstellen.Stringer

type IdLikeGetter interface {
	GetId() schnittstellen.IdLike
}

type DataIdentity interface {
	schnittstellen.ValueLike
	GetTai() ts.Tai
	GetId() IdLike
	schnittstellen.GattungGetter
	GetObjekteSha() schnittstellen.Sha
	GetAkteSha() schnittstellen.Sha
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
	SetTimeAndFields(ts.Tai, ...string) error
	SetFromSku(Sku) error
	SetTransactionIndex(int)
}
