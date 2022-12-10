package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/id_set"
	"github.com/friedenberg/zit/src/kilo/cwd_files"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/zettel_checked_out"
	"github.com/friedenberg/zit/src/oscar/store_fs"
	"github.com/friedenberg/zit/src/papa/umwelt"
	"github.com/friedenberg/zit/src/quebec/user_ops"
)

type Edit struct {
	Or bool
	//TODO-P3 add force
	store_fs.CheckoutMode
}

func init() {
	registerCommand(
		"edit",
		func(f *flag.FlagSet) Command {
			c := &Edit{
				CheckoutMode: store_fs.CheckoutModeZettelOnly,
			}

			f.BoolVar(&c.Or, "or", false, "allow optional criteria instead of required")
			f.Var(&c.CheckoutMode, "mode", "mode for checking out the zettel")

			return commandWithIds{
				CommandWithIds: c,
			}
		},
	)
}

func (c Edit) ProtoIdSet(u *umwelt.Umwelt) (is id_set.ProtoIdSet) {
	is = id_set.MakeProtoIdSet(
		id_set.ProtoId{
			MutableId: &hinweis.Hinweis{},
			Expand: func(v string) (out string, err error) {
				var h hinweis.Hinweis
				h, err = u.StoreObjekten().Abbr().ExpandHinweisString(v)
				out = h.String()
				return
			},
		},
		id_set.ProtoId{
			MutableId: &kennung.Etikett{},
			Expand: func(v string) (out string, err error) {
				var e kennung.Etikett
				e, err = u.StoreObjekten().Abbr().ExpandEtikettString(v)
				out = e.String()
				return
			},
		},
		id_set.ProtoId{
			MutableId: &kennung.Typ{},
		},
		id_set.ProtoId{
			MutableId: &ts.Time{},
		},
	)

	return
}

func (c Edit) RunWithIds(u *umwelt.Umwelt, ids id_set.Set) (err error) {
	checkoutOptions := store_fs.CheckoutOptions{
		CheckoutMode: c.CheckoutMode,
		Format:       zettel.Text{},
	}

	var checkoutResults zettel_checked_out.MutableSet

	query := zettel.WriterIds{
		Filter: id_set.Filter{
			Set: ids,
			Or:  c.Or,
		},
	}

	if checkoutResults, err = u.StoreWorkingDirectory().Checkout(
		checkoutOptions,
		query.WriteZettelTransacted,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = (user_ops.OpenFiles{}).Run(u, checkoutResults.ToSliceFilesAkten()...); err != nil {
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

	files := checkoutResults.ToSliceFilesZettelen()

	if _, err = openVimOp.Run(u, files...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	readOp := user_ops.ReadCheckedOut{
		Umwelt: u,
		OptionsReadExternal: store_fs.OptionsReadExternal{
			Format: zettel.Text{},
		},
	}

	fs := checkoutResults.ToSliceFilesZettelen()

	var possible cwd_files.CwdFiles

	if possible, err = cwd_files.MakeCwdFilesExactly(
		u.Konfig(),
		u.Standort().Cwd(),
		fs...,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	readResults := zettel_checked_out.MakeMutableSetUnique(0)

	if err = readOp.RunMany(possible, readResults.Add); err != nil {
		err = errors.Wrap(err)
		return
	}

	zettels := readResults.ToSliceZettelsExternal()

	checkinOp := user_ops.Checkin{
		Umwelt: u,
		OptionsReadExternal: store_fs.OptionsReadExternal{
			Format: zettel.Text{},
		},
	}

	if _, err = checkinOp.Run(zettels...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
