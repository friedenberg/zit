package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) UpdateTransactedWithExternal(
	kasten kennung.Kasten,
	z *sku.Transacted,
) (err error) {
	switch kasten.GetKastenString() {
	case "chrome":
		err = todo.Implement()

	default:
		if err = s.GetCwdFiles().UpdateTransacted(z); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) ReadTransactedFromKennung(
	k1 schnittstellen.StringerGattungGetter,
) (sk1 *sku.Transacted, err error) {
	sk1 = sku.GetTransactedPool().Get()

	if err = s.ReadOneInto(k1, sk1); err != nil {
		if collections.IsErrNotFound(err) {
			sku.GetTransactedPool().Put(sk1)
			sk1 = nil
		}

		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadTransactedFromKennungKastenSigil(
	k1 schnittstellen.StringerGattungGetter,
	ka kennung.Kasten,
	si kennung.Sigil,
) (sk1 *sku.Transacted, err error) {
	sk1 = sku.GetTransactedPool().Get()

	if err = s.ReadOneInto(k1, sk1); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !si.IncludesExternal() {
		return
	}

	var k3 *kennung.Kennung3

	if k3, err = kennung.MakeKennung3(k1, ka); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ze sku.ExternalLike

	if ze, err = s.ReadOneKennungExternal(
		ObjekteOptions{
			Mode: objekte_mode.ModeUpdateTai,
		},
		k3,
		sk1,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if ze != nil {
		sku.TransactedResetter.ResetWith(sk1, ze.GetSku())
	}

	return
}

func (s *Store) ReadCheckedOutFromTransacted(
	kasten kennung.Kasten,
	sk *sku.Transacted,
) (co sku.CheckedOutLike, err error) {
	switch kasten.GetKastenString() {
	case "chrome":
		err = todo.Implement()

	default:
		if co, err = s.cwdFiles.ReadCheckedOutFromTransacted(sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
