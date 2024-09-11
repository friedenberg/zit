package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) Open(
	repoId ids.RepoId,
	m checkout_mode.Mode,
	ph interfaces.FuncIter[string],
	zsc sku.CheckedOutLikeSet,
) (err error) {
	es, ok := s.externalStores[repoId]

	if !ok {
		err = errors.Errorf("no repo id with id %q", repoId)
		return
	}

	if err = es.Open(m, ph, zsc); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
