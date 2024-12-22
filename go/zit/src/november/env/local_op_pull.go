package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

func (local *Local) PullQueryGroupFromRemote(
	remote Env,
	qg *query.Group,
	printCopies bool,
) (err error) {
	var list *sku.List

	if list, err = remote.MakeInventoryList(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	switch remoteTyped := remote.(type) {
	case *Local:
		importer := local.MakeImporter(remoteTyped.GetDirectoryLayout(), printCopies)

		if err = local.ImportListFromRemoteBlobStore(
			list,
			importer,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = todo.Implement()
		return
	}

	return
}
