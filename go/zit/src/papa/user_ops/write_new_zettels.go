package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type WriteNewZettels struct {
	*umwelt.Umwelt
}

func (c WriteNewZettels) RunMany(
	z zettel.ProtoZettel,
	count int,
) (results sku.TransactedMutableSet, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

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

	return
}

func (c WriteNewZettels) RunOne(
	z zettel.ProtoZettel,
) (result *sku.Transacted, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	if result, err = c.runOneAlreadyLocked(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c WriteNewZettels) runOneAlreadyLocked(
	pz zettel.ProtoZettel,
) (zt *sku.Transacted, err error) {
	zt = pz.Make()

	if err = c.GetStore().CreateOrUpdateFromTransacted(
		zt,
		objekte_mode.ModeApplyProto,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
