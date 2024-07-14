package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/bestandsaufnahme"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Last struct {
	Kasten   ids.RepoId
	Edit     bool
	Organize bool
	Format   string
}

func init() {
	registerCommand(
		"last",
		func(f *flag.FlagSet) Command {
			c := &Last{}

			f.Var(&c.Kasten, "kasten", "none or Chrome")
			f.StringVar(&c.Format, "format", "log", "format")
			f.BoolVar(&c.Organize, "organize", false, "")
			f.BoolVar(&c.Edit, "edit", false, "")

			return c
		},
	)
}

func (c Last) CompletionGattung() ids.Genre {
	return ids.MakeGenre(
		genres.InventoryList,
	)
}

func (c Last) Run(u *umwelt.Umwelt, args ...string) (err error) {
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

	f = iter.MakeSyncSerializer(f)

	if err = c.runWithBestandsaufnahme(u, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.Organize {
		opOrganize := user_ops.Organize{
			Umwelt: u,
		}

		if err = opOrganize.Run(nil, skus); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else if c.Edit {
		opCheckout := user_ops.Checkout{
			Options: checkout_options.Options{
				CheckoutMode: checkout_mode.ModeObjekteAndAkte,
			},
			Umwelt: u,
			Edit:   true,
		}

		if _, err = opCheckout.Run(skus); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c Last) runWithBestandsaufnahme(
	u *umwelt.Umwelt,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	s := u.GetStore()

	var b *sku.Transacted

	if b, err = s.GetBestandsaufnahmeStore().ReadLast(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var a *bestandsaufnahme.InventoryList

	if a, err = s.GetBestandsaufnahmeStore().GetBlob(b.GetBlobSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.TodoP3("support log line format for skus")
	if err = a.Skus.EachPtr(
		func(sk *sku.Transacted) (err error) {
			return f(sk)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
