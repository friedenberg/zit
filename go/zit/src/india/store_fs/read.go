package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) MakeApplyCheckedOut(
	qg sku.Queryable,
	f schnittstellen.FuncIter[sku.CheckedOutLike],
	o sku.ObjekteOptions,
) schnittstellen.FuncIter[*KennungFDPair] {
	return func(em *KennungFDPair) (err error) {
		if err = s.ApplyCheckedOut(o, qg, em, f); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *Store) ApplyCheckedOut(
	o sku.ObjekteOptions,
	qg sku.Queryable,
	em *KennungFDPair,
	f schnittstellen.FuncIter[sku.CheckedOutLike],
) (err error) {
	var co *CheckedOut

	if co, err = s.ReadCheckedOutFromKennungFDPair(o, em); err != nil {
		err = errors.Wrapf(err, "%v", em)
		return
	}

	if !qg.ContainsSku(&co.External.Transacted) {
		return
	}

	if err = f(co); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryCheckedOut(
	qg sku.ExternalQuery,
	f schnittstellen.FuncIter[sku.CheckedOutLike],
) (err error) {
	o := sku.ObjekteOptions{
		Mode: objekte_mode.ModeRealizeSansProto,
	}

	if err = s.All(
		s.MakeApplyCheckedOut(qg, f, o),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryUnsure(
	qg sku.Queryable,
	f schnittstellen.FuncIter[sku.CheckedOutLike],
) (err error) {
	o := sku.ObjekteOptions{
		Mode: objekte_mode.ModeRealizeWithProto,
	}

	if err = s.AllUnsure(
		s.MakeApplyCheckedOut(qg, f, o),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
