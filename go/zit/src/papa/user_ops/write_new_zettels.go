package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/india/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/read_write_repo_local"
)

type WriteNewZettels struct {
	*read_write_repo_local.Repo
}

func (c WriteNewZettels) RunMany(
	z sku.Proto,
	count int,
) (results sku.TransactedMutableSet, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	results = sku.MakeTransactedMutableSet()

	// TODO-P4 modify this to be run once
	for i := 0; i < count; i++ {
		var zt *sku.Transacted

		if zt, err = c.runOneAlreadyLocked(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = results.Add(zt); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = c.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c WriteNewZettels) RunOne(
	z sku.Proto,
) (result *sku.Transacted, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if result, err = c.runOneAlreadyLocked(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c WriteNewZettels) runOneAlreadyLocked(
	pz sku.Proto,
) (zt *sku.Transacted, err error) {
	zt = pz.Make()

	if err = c.GetStore().CreateOrUpdate(
		zt,
		sku.StoreOptions{
			ApplyProto: true,
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
