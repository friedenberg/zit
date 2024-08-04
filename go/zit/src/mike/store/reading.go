package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) UpdateTransactedWithExternal(
	kasten ids.RepoId,
	z *sku.Transacted,
) (err error) {
	kid := kasten.GetRepoIdString()
	es, ok := s.externalStores[kid]

	if !ok {
		err = errors.Errorf("no kasten with id %q", kid)
		return
	}

	if err = es.UpdateTransacted(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadTransactedFromObjectId(
	k1 interfaces.ObjectId,
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

func (s *Store) ReadTransactedFromObjectIdRepoIdSigil(
	k1 interfaces.ObjectId,
	ka ids.RepoId,
	si ids.Sigil,
) (sk1 *sku.Transacted, err error) {
	sk1 = sku.GetTransactedPool().Get()

	if err = s.ReadOneInto(k1, sk1); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !si.IncludesExternal() {
		return
	}

	var ze sku.ExternalLike

	if ze, err = s.ReadOneObjectIdExternal(
		sku.CommitOptions{
			Mode: objekte_mode.ModeUpdateTai,
		},
		k1,
		ka,
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
	kasten ids.RepoId,
	sk *sku.Transacted,
) (co sku.CheckedOutLike, err error) {
	switch kasten.GetRepoIdString() {
	case "browser":
		err = todo.Implement()

	default:
		if co, err = s.cwdFiles.ReadCheckedOutFromTransacted(sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
