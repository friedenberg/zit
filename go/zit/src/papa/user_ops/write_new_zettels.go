package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type WriteNewZettels struct {
	*local_working_copy.Repo
}

func (op WriteNewZettels) RunMany(
	proto sku.Proto,
	count int,
) (results sku.TransactedMutableSet, err error) {
	if err = op.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	results = sku.MakeTransactedMutableSet()

	// TODO-P4 modify this to be run once
	for range count {
		var zt *sku.Transacted

		if zt, err = op.runOneAlreadyLocked(proto); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = results.Add(zt); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = op.Unlock(); err != nil {
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
	proto sku.Proto,
) (object *sku.Transacted, err error) {
	object = proto.Make()

	if err = c.GetStore().CreateOrUpdateDefaultProto(
		object,
		sku.StoreOptions{
			ApplyProto: true,
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
