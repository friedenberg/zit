package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/int_value"
	"github.com/friedenberg/zit/src/echo/ts"
)

type Mutter [2]ts.Time

type IdLike = fmt.Stringer

type DataIdentity interface {
	GetTime() ts.Time
	GetId() IdLike
	schnittstellen.GattungGetter
	GetObjekteSha() schnittstellen.Sha
	GetAkteSha() schnittstellen.Sha
	//TODO-P1 add GetAkteSha
}

type SkuLike interface {
	DataIdentity

	GetKey() string

	GetMutter() Mutter
	GetTransactionIndex() int_value.IntValue
	GetKopf() ts.Time
	GetSchwanz() ts.Time
}

type SkuLikePtr interface {
	SkuLike
	SetTimeAndFields(ts.Time, ...string) error
	SetFromSku(Sku) error
	SetTransactionIndex(int)
}
