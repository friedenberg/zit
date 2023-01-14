package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/int_value"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/ts"
)

type Mutter [2]ts.Time

type IdLike = fmt.Stringer

type DataIdentity interface {
	GetTime() ts.Time
	GetId() IdLike
	gattung.GattungLike
	GetObjekteSha() sha.Sha
	GetAkteSha() sha.Sha
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
