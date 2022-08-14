package user_ops

import (
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/echo/umwelt"
	checkout_store "github.com/friedenberg/zit/golf/store_checkout"
	"github.com/friedenberg/zit/india/store_with_lock"
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
) (result checkout_store.CwdFiles, err error) {
	if result, err = store.CheckoutStore().GetPossibleZettels(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
