package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
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

func (c *PullBlobStore) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c *PullBlobStore) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
	// return ids.MakeGenre(genres.TrueGenre()...)
}

func (cmd *PullBlobStore) Run(
	dep command.Request,
) {
	localWorkingCopy, queryGroup := cmd.MakeLocalWorkingCopyAndQueryGroup(
		dep,
		query.MakeBuilderOptions(cmd),
	)

	importerOptions := store.ImporterOptions{
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

	importer := localWorkingCopy.MakeImporter(importerOptions)

	if err := localWorkingCopy.GetStore().QueryTransacted(
		queryGroup,
		func(sk *sku.Transacted) (err error) {
			if err = importer.ImportBlobIfNecessary(sk); err != nil {
				if dir_layout.IsErrBlobMissing(err) {
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
