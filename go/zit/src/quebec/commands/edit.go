package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/src/charlie/files"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/query"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
	"code.linenisgreat.com/zit/src/papa/user_ops"
)

type Edit struct {
	// TODO-P3 add force
	Delete       bool
	CheckoutMode checkout_mode.Mode
}

func init() {
	registerCommandWithQuery(
		"edit",
		func(f *flag.FlagSet) CommandWithQuery {
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

func (c Edit) CompletionGattung() kennung.Gattung {
	return kennung.MakeGattung(
		gattung.Etikett,
		gattung.Zettel,
		gattung.Typ,
		gattung.Kasten,
	)
}

func (c Edit) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung(
		gattung.Etikett,
		gattung.Zettel,
		gattung.Typ,
		gattung.Kasten,
	)
}

func (c Edit) RunWithQuery(
	u *umwelt.Umwelt,
	ms *query.Group,
) (err error) {
	options := checkout_options.Options{
		CheckoutMode: c.CheckoutMode,
	}

	akten := fd.MakeMutableSet()
	objekten := fd.MakeMutableSet()

	if err = u.GetStore().CheckoutQuery(
		options,
		ms,
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

	objektenFiles := iter.Strings[*fd.FD](objekten)
	aktenFiles := iter.Strings[*fd.FD](akten)

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
