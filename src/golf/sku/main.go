package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/int_value"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/schnittstellen"
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
	SetTimeAndFields(ts.Time, ...string) error
	SetFromSku(Sku) error

	GetKey() string

	GetMutter() Mutter
	SetTransactionIndex(int)
	GetTransactionIndex() int_value.IntValue
	GetKopf() ts.Time
	GetSchwanz() ts.Time
}
