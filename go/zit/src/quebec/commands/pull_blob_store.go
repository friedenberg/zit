package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("pull-blob-store", &PullBlobStore{})
}

type PullBlobStore struct {
	command_components.LocalWorkingCopyWithQueryGroup
	command_components.RemoteBlobStore
}

func (cmd *PullBlobStore) SetFlagSet(f *flag.FlagSet) {
	cmd.LocalWorkingCopyWithQueryGroup.SetFlagSet(f)
	cmd.RemoteBlobStore.SetFlagSet(f)
}

func (cmd *PullBlobStore) Run(
	dep command.Request,
) {
	localWorkingCopy, queryGroup := cmd.MakeLocalWorkingCopyAndQueryGroup(
		dep,
		query.BuilderOptionsOld(
			cmd,
			query.BuilderOptionDefaultSigil(
				ids.SigilHistory,
				ids.SigilHidden,
			),
			query.BuilderOptionDefaultGenres(genres.InventoryList),
		),
	)

	importerOptions := sku.ImporterOptions{
		ExcludeObjects: true,
		PrintCopies:    true,
	}

	{
		var err error

		if importerOptions.RemoteBlobStore, err = cmd.MakeRemoteBlobStore(
			localWorkingCopy,
		); err != nil {
			dep.CancelWithError(err)
		}
	}

	importer := localWorkingCopy.MakeImporter(
		importerOptions,
		sku.GetStoreOptionsRemoteTransfer(),
	)

	if err := localWorkingCopy.GetStore().QueryTransacted(
		queryGroup,
		func(sk *sku.Transacted) (err error) {
			if err = importer.ImportBlobIfNecessary(sk); err != nil {
				if env_dir.IsErrBlobMissing(err) {
					err = nil
					localWorkingCopy.GetUI().Printf("Blob missing from remote: %q", sk.GetBlobSha())
				} else {
					err = errors.Wrap(err)
				}

				return
			}

			return
		},
	); err != nil {
		dep.CancelWithError(err)
	}
}
