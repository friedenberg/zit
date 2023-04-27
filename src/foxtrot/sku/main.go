package sku

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/echo/ts"
)

type Mutter [2]ts.Time

type IdLike = schnittstellen.Stringer

type IdLikeGetter interface {
	GetId() schnittstellen.IdLike
}

type DataIdentity interface {
	schnittstellen.ValueLike
	GetTime() ts.Time
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

	GetMutter() Mutter
	GetTransactionIndex() values.Int
	GetKopf() ts.Time
	GetSchwanz() ts.Time
}

type SkuLikePtr interface {
	SkuLike
	SetTimeAndFields(ts.Time, ...string) error
	SetFromSku(Sku) error
	SetTransactionIndex(int)
}
