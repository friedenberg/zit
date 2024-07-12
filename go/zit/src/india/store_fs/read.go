package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
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

	// ui.Debug().Print(qg, qg.ContainsSku(&co.External.Transacted), co)

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
	qg *query.Group,
	f schnittstellen.FuncIter[sku.CheckedOutLike],
) (err error) {
	wg := iter.MakeErrorWaitGroupParallel()

	{
		o := sku.ObjekteOptions{
			Mode: objekte_mode.ModeRealizeSansProto,
		}

		wg.Do(func() error {
			return s.All(s.MakeApplyCheckedOut(qg, f, o))
		})
	}

	if !qg.ExcludeUntracked {
		wg.Do(func() error {
			return s.QueryUnsure(qg, f)
		})
	}

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO [cot/gl !task project-2021-zit-kasten today zz-inbox] move unsure akten and untracked into kasten interface and store_fs
func (s *Store) QueryUnsure(
	qg *query.Group,
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
