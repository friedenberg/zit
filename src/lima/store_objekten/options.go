package store_objekten

import (
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/kilo/cwd"
)

type CheckoutOptions struct {
	Cwd          cwd.CwdFiles
	Force        bool
	CheckoutMode objekte.CheckoutMode
}
