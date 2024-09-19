package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) SaveBlob(el sku.ExternalLike) (err error) {
	repoId := el.GetRepoId()
	es, ok := s.externalStores[repoId]

	if !ok {
		err = errors.Errorf("no kasten with id %q", repoId)
		return
	}

	if err = es.SaveBlob(el); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
