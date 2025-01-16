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
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/lima/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Last struct {
	command_components.RepoLayout
	command_components.LocalArchive

	RepoId   ids.RepoId
	Edit     bool
	Organize bool
	Format   string
}

func init() {
	registerCommand(
		"last",
		&Last{},
	)
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

func (cmd Last) Run(dependencies Dependencies) {
	repoLayout := cmd.MakeRepoLayout(
		dependencies.Context,
		dependencies.Config,
		false,
	)

	archive := cmd.MakeLocalArchive(repoLayout)

	if len(dependencies.Args()) != 0 {
		ui.Err().Print("ignoring arguments")
	}

	if localWorkingCopy, ok := archive.(*local_working_copy.Repo); ok {
		cmd.runLocalWorkingCopy(localWorkingCopy)
	} else {
		cmd.runArchive(archive)
	}
}

func (c Last) runArchive(archive repo.Archive) {
	if (c.Edit || c.Organize) && c.Format != "" {
		archive.GetRepoLayout().CancelWithErrorf("cannot organize, edit, or specify format for Archive repos")
	}

	boxFormat := box_format.MakeBoxTransactedArchive(
		archive.GetRepoLayout().Env,
		options_print.V0{}.WithPrintTai(true),
	)

	f := string_format_writer.MakeDelim(
		"\n",
		archive.GetRepoLayout().GetUIFile(),
		string_format_writer.MakeFunc(
			func(w interfaces.WriterAndStringWriter, o *sku.Transacted) (n int64, err error) {
				return boxFormat.WriteStringFormat(w, o)
			},
		),
	)

	f = quiter.MakeSyncSerializer(f)

	if err := c.runWithInventoryList(archive, f); err != nil {
		archive.GetRepoLayout().CancelWithError(err)
	}
}

func (c Last) runLocalWorkingCopy(archive *local_working_copy.Repo) {
	if (c.Edit || c.Organize) && c.Format != "" {
		ui.Err().Print("ignoring format")
	} else if c.Edit && c.Organize {
		archive.GetRepoLayout().CancelWithErrorf("cannot organize and edit at the same time")
	}

	skus := sku.MakeTransactedMutableSet()

	var f interfaces.FuncIter[*sku.Transacted]

	if c.Organize || c.Edit {
		f = skus.Add
	} else {
		{
			var err error

			if f, err = archive.MakeFormatFunc(
				c.Format,
				archive.GetRepoLayout().GetUIFile(),
			); err != nil {
				archive.GetRepoLayout().CancelWithError(err)
			}
		}
	}

	f = quiter.MakeSyncSerializer(f)

	if err := c.runWithInventoryList(archive, f); err != nil {
		archive.GetRepoLayout().CancelWithError(err)
	}

	if c.Organize {
		opOrganize := user_ops.Organize{
			Repo: archive,
			Metadata: organize_text.Metadata{
				OptionCommentSet: organize_text.MakeOptionCommentSet(nil),
			},
		}

		var results organize_text.OrganizeResults

		{
			var err error

			if results, err = opOrganize.RunWithTransacted(nil, skus); err != nil {
				archive.GetRepoLayout().CancelWithError(err)
			}
		}

		if _, err := archive.LockAndCommitOrganizeResults(results); err != nil {
			archive.GetRepoLayout().CancelWithError(err)
		}
	} else if c.Edit {
		opCheckout := user_ops.Checkout{
			Options: checkout_options.Options{
				CheckoutMode: checkout_mode.MetadataAndBlob,
			},
			Repo: archive,
			Edit: true,
		}

		if _, err := opCheckout.Run(skus); err != nil {
			archive.GetRepoLayout().CancelWithError(err)
		}
	}
}

func (c Last) runWithInventoryList(
	archive repo.Archive,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var b *sku.Transacted

	inventoryListStore := archive.GetInventoryListStore()

	if b, err = inventoryListStore.ReadLast(); err != nil {
		err = errors.Wrap(err)
		return
	}

	objectFormat := object_inventory_format.FormatForVersion(
		archive.GetRepoLayout().GetStoreVersion(),
	)

	boxFormat := box_format.MakeBoxTransactedArchive(
		archive.GetRepoLayout().Env,
		options_print.V0{}.WithPrintTai(true),
	)

	inventoryListBlobStore := blob_store.MakeInventoryStore(
		archive.GetRepoLayout(),
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
