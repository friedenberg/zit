package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
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

	importer := local.MakeImporter(printCopies)

	importer.RemoteBlobStore = remote.GetBlobStore()
	importer.ParentNegotiator = ParentNegotiatorFirstAncestor{
		Local:  local,
		Remote: remote,
	}

	if err = local.ImportList(
		list,
		importer,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
