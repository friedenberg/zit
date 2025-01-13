package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/read_write_repo_local"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Last struct {
	RepoId   ids.RepoId
	Edit     bool
	Organize bool
	Format   string
}

func init() {
	registerCommand(
		"last",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &Last{}

			f.Var(&c.RepoId, "kasten", "none or Browser")
			f.StringVar(&c.Format, "format", "log", "format")
			f.BoolVar(&c.Organize, "organize", false, "")
			f.BoolVar(&c.Edit, "edit", false, "")

			return c
		},
	)
}

func (c Last) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.InventoryList,
	)
}

func (c Last) RunWithRepo(u *read_write_repo_local.Repo, args ...string) {
	if len(args) != 0 {
		ui.Err().Print("ignoring arguments")
	}

	if (c.Edit || c.Organize) && c.Format != "" {
		ui.Err().Print("ignoring format")
	} else if c.Edit && c.Organize {
		u.CancelWithErrorf("cannot organize and edit at the same time")
	}

	skus := sku.MakeTransactedMutableSet()

	var f interfaces.FuncIter[*sku.Transacted]

	if c.Organize || c.Edit {
		f = skus.Add
	} else {
		{
			var err error

			if f, err = u.MakeFormatFunc(c.Format, u.GetUIFile()); err != nil {
				u.CancelWithError(err)
			}
		}
	}

	f = quiter.MakeSyncSerializer(f)

	if err := c.runWithInventoryList(u, f); err != nil {
		u.CancelWithError(err)
	}

	if c.Organize {
		opOrganize := user_ops.Organize{
			Repo: u,
			Metadata: organize_text.Metadata{
				OptionCommentSet: organize_text.MakeOptionCommentSet(nil),
			},
		}

		var results organize_text.OrganizeResults

		{
			var err error

			if results, err = opOrganize.RunWithTransacted(nil, skus); err != nil {
				u.CancelWithError(err)
			}
		}

		if _, err := u.LockAndCommitOrganizeResults(results); err != nil {
			u.CancelWithError(err)
		}
	} else if c.Edit {
		opCheckout := user_ops.Checkout{
			Options: checkout_options.Options{
				CheckoutMode: checkout_mode.MetadataAndBlob,
			},
			Repo: u,
			Edit: true,
		}

		if _, err := opCheckout.Run(skus); err != nil {
			u.CancelWithError(err)
		}
	}
}

func (c Last) runWithInventoryList(
	u *read_write_repo_local.Repo,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	s := u.GetStore()

	var b *sku.Transacted

	if b, err = s.GetInventoryListStore().ReadLast(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var twb sku.TransactedWithBlob[*sku.List]

	if twb, _, err = s.GetBlobStore().GetInventoryList().GetTransactedWithBlob(
		b,
	); err != nil {
		err = errors.Wrapf(err, "InventoryList: %q", b)
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
