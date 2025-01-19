package local_working_copy

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/mike/importer"
)

func (repo *Repo) ImportList(
	list *sku.List,
	i sku.Importer,
) (err error) {
	repo.Must(repo.Lock)

	if err = repo.GetInventoryListStore().ImportList(
		list,
		i,
	); err != nil {
		if !errors.Is(err, importer.ErrNeedsMerge) {
			err = errors.Wrap(err)
			return
		}
	}

	repo.Must(repo.Unlock)

	return
}
