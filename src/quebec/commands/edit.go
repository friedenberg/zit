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
	"github.com/friedenberg/zit/src/juliett/zettel"
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
	registerCommandWithQuery(
		"edit",
		func(f *flag.FlagSet) CommandWithQuery {
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

func (c Edit) RunWithQuery(u *umwelt.Umwelt, ms kennung.MetaSet) (err error) {
	ids, ok := ms.Get(gattung.Zettel)

	if !ok {
		return
	}

	return c.editZettels(u, ids)
}

func (c Edit) runWithQuery(u *umwelt.Umwelt, ms kennung.MetaSet) (err error) {
	checkoutOptions := store_fs.CheckoutOptions{
		CheckoutMode: c.CheckoutMode,
	}

	akten := kennung.MakeMutableFDSet()
	objekten := kennung.MakeMutableFDSet()

	if err = u.StoreWorkingDirectory().CheckoutQuery(
		checkoutOptions,
		ms,
		func(co objekte.CheckedOutLike) (err error) {
			e := co.GetExternal()

			akten.Add(e.GetAkteFD())
			objekten.Add(e.GetObjekteFD())

			return
		},
	); err != nil {
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

	// readOp := user_ops.ReadCheckedOut{
	// 	Umwelt:              u,
	// 	OptionsReadExternal: store_fs.OptionsReadExternal{},
	// }

	// var possible cwd.CwdFiles

	filez := append([]string{}, objektenFiles...)
	filez = append(filez, aktenFiles...)

	if _, err = cwd.MakeCwdFilesExactly(
		u.Konfig(),
		u.Standort().Cwd(),
		filez...,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Edit) editZettels(u *umwelt.Umwelt, ids kennung.Set) (err error) {
	checkoutOptions := store_fs.CheckoutOptions{
		CheckoutMode: c.CheckoutMode,
	}

	var checkoutResults zettel.MutableSetCheckedOut

	query := zettel.WriterIds{
		Filter: kennung.Filter{
			Set: ids,
		},
	}

	if checkoutResults, err = u.StoreWorkingDirectory().Checkout(
		checkoutOptions,
		query.WriteZettelTransacted,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var filesAkten []string

	if filesAkten, err = zettel.ToSliceFilesAkten(checkoutResults); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = (user_ops.OpenFiles{}).Run(u, filesAkten...); err != nil {
		err = errors.Wrap(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithCursorLocation(2, 3).
			WithFileType("zit-zettel").
			WithInsertMode().
			Build(),
	}

	var filesZettelen []string

	if filesZettelen, err = zettel.ToSliceFilesZettelen(checkoutResults); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = openVimOp.Run(u, filesZettelen...); err != nil {
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

	var possible cwd.CwdFiles

	if possible, err = cwd.MakeCwdFilesExactly(
		u.Konfig(),
		u.Standort().Cwd(),
		filesZettelen...,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	checkinOp := user_ops.Checkin{
		Delete: c.Delete,
	}

	if err = checkinOp.Run(u, possible); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
