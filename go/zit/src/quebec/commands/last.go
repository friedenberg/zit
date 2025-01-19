package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

func init() {
	command.Register("last", &Last{})
}

type Last struct {
	command_components.RepoLayout
	command_components.LocalArchive

	RepoId   ids.RepoId
	Edit     bool
	Organize bool
	Format   string
}

func (cmd *Last) SetFlagSet(f *flag.FlagSet) {
	cmd.RepoLayout.SetFlagSet(f)
	cmd.LocalArchive.SetFlagSet(f)

	f.Var(&cmd.RepoId, "kasten", "none or Browser")
	f.StringVar(&cmd.Format, "format", "log", "format")
	f.BoolVar(&cmd.Organize, "organize", false, "")
	f.BoolVar(&cmd.Edit, "edit", false, "")
}

func (c Last) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.InventoryList,
	)
}

func (cmd Last) Run(dep command.Request) {
	repoLayout := cmd.MakeRepoLayout(dep, false)

	archive := cmd.MakeLocalArchive(repoLayout)

	if len(dep.Args()) != 0 {
		ui.Err().Print("ignoring arguments")
	}

	if localWorkingCopy, ok := archive.(*local_working_copy.Repo); ok {
		cmd.runLocalWorkingCopy(localWorkingCopy)
	} else {
		cmd.runArchive(repoLayout, archive)
	}
}

func (c Last) runArchive(repoLayout env_repo.Env, archive repo.Repo) {
	if (c.Edit || c.Organize) && c.Format != "" {
		repoLayout.CancelWithErrorf("cannot organize, edit, or specify format for Archive repos")
	}

	boxFormat := box_format.MakeBoxTransactedArchive(
		repoLayout,
		options_print.V0{}.WithPrintTai(true),
	)

	f := string_format_writer.MakeDelim(
		"\n",
		repoLayout.GetUIFile(),
		string_format_writer.MakeFunc(
			func(w interfaces.WriterAndStringWriter, o *sku.Transacted) (n int64, err error) {
				return boxFormat.WriteStringFormat(w, o)
			},
		),
	)

	f = quiter.MakeSyncSerializer(f)

	if err := c.runWithInventoryList(repoLayout, archive, f); err != nil {
		repoLayout.CancelWithError(err)
	}
}

func (c Last) runLocalWorkingCopy(localWorkingCopy *local_working_copy.Repo) {
	if (c.Edit || c.Organize) && c.Format != "" {
		ui.Err().Print("ignoring format")
	} else if c.Edit && c.Organize {
		localWorkingCopy.GetRepoLayout().CancelWithErrorf("cannot organize and edit at the same time")
	}

	skus := sku.MakeTransactedMutableSet()

	var f interfaces.FuncIter[*sku.Transacted]

	if c.Organize || c.Edit {
		f = skus.Add
	} else {
		{
			var err error

			if f, err = localWorkingCopy.MakeFormatFunc(
				c.Format,
				localWorkingCopy.GetRepoLayout().GetUIFile(),
			); err != nil {
				localWorkingCopy.GetRepoLayout().CancelWithError(err)
			}
		}
	}

	f = quiter.MakeSyncSerializer(f)

	if err := c.runWithInventoryList(
		localWorkingCopy.GetRepoLayout(),
		localWorkingCopy,
		f,
	); err != nil {
		localWorkingCopy.GetRepoLayout().CancelWithError(err)
	}

	if c.Organize {
		opOrganize := user_ops.Organize{
			Repo: localWorkingCopy,
			Metadata: organize_text.Metadata{
				OptionCommentSet: organize_text.MakeOptionCommentSet(nil),
			},
		}

		var results organize_text.OrganizeResults

		{
			var err error

			if results, err = opOrganize.RunWithTransacted(nil, skus); err != nil {
				localWorkingCopy.GetRepoLayout().CancelWithError(err)
			}
		}

		if _, err := localWorkingCopy.LockAndCommitOrganizeResults(results); err != nil {
			localWorkingCopy.GetRepoLayout().CancelWithError(err)
		}
	} else if c.Edit {
		opCheckout := user_ops.Checkout{
			Options: checkout_options.Options{
				CheckoutMode: checkout_mode.MetadataAndBlob,
			},
			Repo: localWorkingCopy,
			Edit: true,
		}

		if _, err := opCheckout.Run(skus); err != nil {
			localWorkingCopy.GetRepoLayout().CancelWithError(err)
		}
	}
}

func (c Last) runWithInventoryList(
	repoLayout env_repo.Env,
	archive repo.Repo,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var b *sku.Transacted

	inventoryListStore := archive.GetInventoryListStore()

	if b, err = inventoryListStore.ReadLast(); err != nil {
		err = errors.Wrap(err)
		return
	}

	objectFormat := object_inventory_format.FormatForVersion(
		repoLayout.GetStoreVersion(),
	)

	boxFormat := box_format.MakeBoxTransactedArchive(
		repoLayout,
		options_print.V0{}.WithPrintTai(true),
	)

	inventoryListBlobStore := typed_blob_store.MakeInventoryStore(
		repoLayout,
		objectFormat,
		boxFormat,
	)

	var twb sku.TransactedWithBlob[*sku.List]

	if twb, _, err = inventoryListBlobStore.GetTransactedWithBlob(
		b,
	); err != nil {
		err = errors.Wrapf(err, "InventoryList: %q", sku.String(b))
		return
	}

	ui.TodoP3("support log line format for skus")
	if err = twb.Blob.EachPtr(
		func(sk *sku.Transacted) (err error) {
			return f(sk)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
