package store_objekten

import (
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/lima/cwd"
)

type CheckoutOptions struct {
	Cwd          cwd.CwdFiles
	Force        bool
	CheckoutMode checkout_mode.Mode
}
