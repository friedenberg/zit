package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type ObjekteOptions = sku.ObjectOptions

func (s *Store) ReadOneObjectIdExternal(
	o ObjekteOptions,
	k1 interfaces.ObjectIdWithRepoId,
	sk *sku.Transacted,
) (el sku.ExternalLike, err error) {
	switch k1.GetRepoId().GetRepoIdString() {
	case "chrome":
		// TODO populate with chrome kasten
		ui.Debug().Print("would populate from chrome")

	default:
		if el, err = s.cwdFiles.ReadTransactedFromObjectId(o, k1, sk); err != nil {
			err = errors.Wrap(err)
			return
		}
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
