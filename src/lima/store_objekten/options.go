package store_objekten

import (
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/kilo/cwd"
)

type CheckoutOptions struct {
	Cwd          cwd.CwdFiles
	Force        bool
	CheckoutMode sku.CheckoutMode
}
