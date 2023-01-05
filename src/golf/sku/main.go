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

type SkuLike interface {
	gattung.GattungLike
	SetFields(...string) error
	GetKey() string
	SetTransactionIndex(int)
	GetId() IdLike
	GetMutter() Mutter
	GetObjekteSha() sha.Sha
	GetTransactionIndex() int_value.IntValue
	GetKopf() ts.Time
	GetSchwanz() ts.Time
}

type FuncSkuObjekteReader func(SkuLike) (sha.ReadCloser, error)
