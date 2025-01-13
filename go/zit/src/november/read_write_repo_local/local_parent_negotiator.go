package read_write_repo_local

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/india/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
)

type ParentNegotiatorFirstAncestor struct {
	Local, Remote repo.ReadWrite
}

func (parentNegotiator ParentNegotiatorFirstAncestor) GetParentNegotiator() sku.ParentNegotiator {
	return parentNegotiator
}

func (parentNegotiator ParentNegotiatorFirstAncestor) FindBestCommonAncestor(
	conflicted sku.Conflicted,
) (ancestor *sku.Transacted, err error) {
	var ancestorsLocal, ancestorsRemote []*sku.Transacted

	wg := errors.MakeWaitGroupParallel()

	wg.Do(
		func() (err error) {
			if ancestorsLocal, err = parentNegotiator.Local.ReadObjectHistory(
				conflicted.Local.GetObjectId(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)

	wg.Do(
		func() (err error) {
			if ancestorsRemote, err = parentNegotiator.Remote.ReadObjectHistory(
				conflicted.Local.GetObjectId(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if len(ancestorsLocal) == 0 || len(ancestorsRemote) == 0 {
		return
	}

	// TODO repool all skus except ancestor

	ancestorLocal := ancestorsLocal[len(ancestorsLocal)-1]
	ancestorRemote := ancestorsRemote[len(ancestorsRemote)-1]

	if ancestorLocal.Metadata.EqualsSansTai(&ancestorRemote.Metadata) {
		ancestor = ancestorLocal
	}

	return
}
