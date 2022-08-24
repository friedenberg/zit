package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/echo/umwelt"
	store_checkout "github.com/friedenberg/zit/src/hotel/store_checkout"
	"github.com/friedenberg/zit/src/india/store_with_lock"
)

type GetPossibleZettels struct {
	umwelt *umwelt.Umwelt
}

func NewGetPossibleZettels(u *umwelt.Umwelt) GetPossibleZettels {
	return GetPossibleZettels{
		umwelt: u,
	}
}

func (uo GetPossibleZettels) Run(
	store store_with_lock.Store,
) (result store_checkout.CwdFiles, err error) {
	if result, err = store.CheckoutStore().GetPossibleZettels(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
