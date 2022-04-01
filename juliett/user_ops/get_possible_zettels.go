package user_ops

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/delta/umwelt"
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

func (uo GetPossibleZettels) Run() (hinweisen []string, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(uo.umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	hinweisen, err = store.CheckoutStore().GetPossibleZettels()

	return
}
