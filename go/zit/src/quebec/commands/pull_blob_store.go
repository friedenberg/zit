package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type PullBlobStore struct {
	command_components.RemoteBlobStore
}

func init() {
	registerCommandWithQuery(
		"pull-blob-store",
		func(f *flag.FlagSet) CommandWithQuery {
			cmd := &PullBlobStore{}
			cmd.SetFlagSet(f)
			return cmd
		},
	)
}

func (cmd *PullBlobStore) SetFlagSet(f *flag.FlagSet) {
	cmd.RemoteBlobStore.SetFlagSet(f)
}

func (c *PullBlobStore) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c *PullBlobStore) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
	// return ids.MakeGenre(genres.TrueGenre()...)
}

func (c *PullBlobStore) RunWithQuery(
	repo *repo_local_working_copy.Repo,
	qg *query.Group,
) {
	importerOptions := store.ImporterOptions{
		ExcludeObjects: true,
		PrintCopies:    true,
	}

	{
		var err error

		if importerOptions.RemoteBlobStore, err = c.MakeRemoteBlobStore(
			repo.Env,
		); err != nil {
			repo.CancelWithError(err)
		}
	}

	importer := repo.MakeImporter(importerOptions)

	if err := repo.GetStore().QueryTransacted(
		qg,
		func(sk *sku.Transacted) (err error) {
			if err = importer.ImportBlobIfNecessary(sk); err != nil {
				if dir_layout.IsErrBlobMissing(err) {
					err = nil
					repo.GetUI().Printf("Blob missing from remote: %q", sk.GetBlobSha())
				} else {
					err = errors.Wrap(err)
				}

				return
			}

			return
		},
	); err != nil {
		repo.CancelWithError(err)
	}
}
