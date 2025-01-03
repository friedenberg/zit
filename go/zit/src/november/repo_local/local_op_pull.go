package repo_local

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (local *Repo) PullQueryGroupFromRemote(
	remote repo.Repo,
	qg *query.Group,
	options repo.RemoteTransferOptions,
) (err error) {
	var list *sku.List

	if list, err = remote.MakeInventoryList(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	importerOptions := store.ImporterOptions{
		CheckedOutPrinter:   local.PrinterCheckedOutConflictsForRemoteTransfers(),
		AllowMergeConflicts: options.AllowMergeConflicts,
		ExcludeObjects:      !options.IncludeObjects,
	}

	if options.IncludeBlobs {
		importerOptions.RemoteBlobStore = remote.GetBlobStore()
	}

	importerOptions.ParentNegotiator = ParentNegotiatorFirstAncestor{
		Local:  local,
		Remote: remote,
	}

	importerOptions.PrintCopies = options.PrintCopies
	importer := local.MakeImporter(importerOptions)

	if err = local.ImportList(
		list,
		importer,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
