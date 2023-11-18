package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/checkout_options"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/oscar/umwelt"
	"github.com/friedenberg/zit/src/papa/user_ops"
)

type Edit struct {
	// TODO-P3 add force
	Delete       bool
	CheckoutMode checkout_mode.Mode
}

func init() {
	registerCommandWithCwdQuery(
		"edit",
		func(f *flag.FlagSet) CommandWithCwdQuery {
			c := &Edit{
				CheckoutMode: checkout_mode.ModeObjekteOnly,
			}

			f.BoolVar(
				&c.Delete,
				"delete",
				false,
				"delete the zettel and akte after successful checkin",
			)
			f.Var(&c.CheckoutMode, "mode", "mode for checking out the zettel")

			return c
		},
	)
}

func (c Edit) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Etikett,
		gattung.Zettel,
		gattung.Typ,
		gattung.Kasten,
	)
}

func (c Edit) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Etikett,
		gattung.Zettel,
		gattung.Typ,
		gattung.Kasten,
	)
}

func (c Edit) RunWithCwdQuery(
	u *umwelt.Umwelt,
	ms matcher.Query,
	pz *cwd.CwdFiles,
) (err error) {
	options := checkout_options.Options{
		CheckoutMode: c.CheckoutMode,
	}

	akten := fd.MakeMutableSet()
	objekten := fd.MakeMutableSet()

	if err = u.StoreObjekten().CheckoutQuery(
		options,
		matcher.MakeFuncReaderTransactedLikePtr(ms, u.StoreObjekten().QueryWithCwd),
		func(co *sku.CheckedOut) (err error) {
			e := co.External

			if afd := e.GetAkteFD(); afd.String() != "." {
				akten.Add(afd)
			}

			if ofd := e.GetObjekteFD(); ofd.String() != "." {
				objekten.Add(ofd)
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	objektenFiles := iter.Strings[fd.FD](objekten)
	aktenFiles := iter.Strings[fd.FD](akten)

	if err = (user_ops.OpenFiles{}).Run(u, aktenFiles...); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithCursorLocation(2, 3).
			WithInsertMode().
			Build(),
	}

	if _, err = openVimOp.Run(u, objektenFiles...); err != nil {
		if errors.Is(err, files.ErrEmptyFileList) {
			err = errors.Normalf("nothing to open in vim")
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	if err = u.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	filez := append([]string{}, objektenFiles...)
	filez = append(filez, aktenFiles...)

	op := user_ops.Checkin{
		Delete: c.Delete,
	}

	if err = op.Run(u, ms); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
