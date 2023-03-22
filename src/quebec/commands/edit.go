package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/oscar/user_ops"
)

type Edit struct {
	// TODO-P3 add force
	Delete       bool
	CheckoutMode objekte.CheckoutMode
}

func init() {
	registerCommandWithCwdQuery(
		"edit",
		func(f *flag.FlagSet) CommandWithCwdQuery {
			c := &Edit{
				CheckoutMode: objekte.CheckoutModeObjekteOnly,
			}

			f.BoolVar(&c.Delete, "delete", false, "delete the zettel and akte after successful checkin")
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
		gattung.Konfig,
	)
}

func (c Edit) RunWithCwdQuery(
	u *umwelt.Umwelt,
	ms kennung.MetaSet,
	pz cwd.CwdFiles,
) (err error) {
	options := store_fs.CheckoutOptions{
		Cwd:          pz,
		CheckoutMode: c.CheckoutMode,
	}

	akten := kennung.MakeMutableFDSet()
	objekten := kennung.MakeMutableFDSet()

	if err = u.StoreWorkingDirectory().CheckoutQuery(
		options,
		ms,
		func(co objekte.CheckedOutLike) (err error) {
			e := co.GetExternal()

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

	objektenFiles := collections.Strings[kennung.FD](objekten)
	aktenFiles := collections.Strings[kennung.FD](akten)

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

	var cwdFiles cwd.CwdFiles

	if cwdFiles, err = cwd.MakeCwdFilesExactly(
		u.Konfig(),
		u.Standort().Cwd(),
		filez...,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	op := user_ops.Checkin{
		Delete: c.Delete,
	}

	if err = op.Run(u, ms, cwdFiles); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
