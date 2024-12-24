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
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
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
		func(f *flag.FlagSet) Command {
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

func (c Last) Run(u *env.Local, args ...string) (err error) {
	if len(args) != 0 {
		ui.Err().Print("ignoring arguments")
	}

	if (c.Edit || c.Organize) && c.Format != "" {
		ui.Err().Print("ignoring format")
	} else if c.Edit && c.Organize {
		err = errors.Errorf("cannot organize and edit at the same time")
		return
	}

	skus := sku.MakeTransactedMutableSet()

	var f interfaces.FuncIter[*sku.Transacted]

	if c.Organize || c.Edit {
		f = skus.Add
	} else {
		if f, err = u.MakeFormatFunc(c.Format, u.Out()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	f = quiter.MakeSyncSerializer(f)

	if err = c.runWithInventoryList(u, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.Organize {
		opOrganize := user_ops.Organize{
			Local: u,
			Metadata: organize_text.Metadata{
				OptionCommentSet: organize_text.MakeOptionCommentSet(nil),
			},
		}

		var results organize_text.OrganizeResults

		if results, err = opOrganize.RunWithTransacted(nil, skus); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = u.LockAndCommitOrganizeResults(results); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else if c.Edit {
		opCheckout := user_ops.Checkout{
			Options: checkout_options.Options{
				CheckoutMode: checkout_mode.MetadataAndBlob,
			},
			Local:  u,
			Edit: true,
		}

		if _, err = opCheckout.Run(skus); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c Last) runWithInventoryList(
	u *env.Local,
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
