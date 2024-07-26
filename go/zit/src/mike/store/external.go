package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) ReadOneObjectIdExternal(
	o sku.CommitOptions,
	oid interfaces.ObjectId,
	kid ids.RepoId,
	sk *sku.Transacted,
) (el sku.ExternalLike, err error) {
	es, ok := s.externalStores[kid.GetRepoIdString()]

	if !ok {
		err = errors.Errorf("no kasten with id %q", kid)
		return
	}

	if el, err = es.ReadTransactedFromObjectId(o, oid, sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) Open(
	kasten interfaces.RepoIdGetter,
	m checkout_mode.Mode,
	ph interfaces.FuncIter[string],
	zsc sku.CheckedOutLikeSet,
) (err error) {
	kid := kasten.GetRepoId().GetRepoIdString()
	es, ok := s.externalStores[kid]

	if !ok {
		err = errors.Errorf("no kasten with id %q", kid)
		return
	}

	if err = es.Open(m, ph, zsc); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
