package chrome

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) ReadIntoExternalFromKennungFDPair(
	o sku.ObjekteOptions,
	em *item,
	t *sku.Transacted,
	e *External,
) (err error) {
	o.Del(objekte_mode.ModeApplyProto)

	if err = s.ReadOneExternalInto(
		&o,
		em,
		t,
		e,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.externalStoreInfo.FuncCommit(
		&e.Transacted,
		o,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOneExternalInto(
	o *sku.ObjekteOptions,
	em *item,
	t *sku.Transacted,
	e *External,
) (err error) {
	if o.Clock == nil {
	}

	return
}
