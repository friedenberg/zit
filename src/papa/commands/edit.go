package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/konfig"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/typ"
	"github.com/friedenberg/zit/src/foxtrot/id_set"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/hotel/zettel_named"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	store_fs "github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
	"github.com/friedenberg/zit/src/oscar/user_ops"
)

type Edit struct {
	Or bool
	//TODO add force
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
			MutableId: &konfig.Id{},
		},
		id_set.ProtoId{
			MutableId: &hinweis.Hinweis{},
			Expand: func(v string) (out string, err error) {
				var h hinweis.Hinweis
				h, err = u.StoreObjekten().ExpandHinweisString(v)
				out = h.String()
				return
			},
		},
		id_set.ProtoId{
			MutableId: &etikett.Etikett{},
			Expand: func(v string) (out string, err error) {
				var e etikett.Etikett
				e, err = u.StoreObjekten().ExpandEtikettString(v)
				out = e.String()
				return
			},
		},
		id_set.ProtoId{
			MutableId: &typ.Typ{},
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

	query := zettel_transacted.WriterIds(
		zettel_named.FilterIdSet{
			Set: ids,
			Or:  c.Or,
		},
	)

	if checkoutResults, err = u.StoreWorkingDirectory().Checkout(
		checkoutOptions,
		query.WriteZettelTransacted,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = (user_ops.OpenFiles{}).Run(checkoutResults.ToSliceFilesAkten()...); err != nil {
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

	if ids.HasKonfig() {
		files = append(files, u.Standort().FileKonfigToml())
	}

	if _, err = openVimOp.Run(files...); err != nil {
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

	var possible store_fs.CwdFiles

	if possible, err = store_fs.MakeCwdFilesExactly(u.Konfig().Compiled, u.Standort().Cwd(), fs...); err != nil {
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