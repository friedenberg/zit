package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/charlie/typ"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/golf/zettel_external"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/india/store_working_directory"
	"github.com/friedenberg/zit/src/kilo/umwelt"
	"github.com/friedenberg/zit/src/lima/user_ops"
)

type Edit struct {
	Or bool
	//TODO add force
	store_working_directory.CheckoutMode
}

func init() {
	registerCommand(
		"edit",
		func(f *flag.FlagSet) Command {
			c := &Edit{
				CheckoutMode: store_working_directory.CheckoutModeZettelOnly,
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
	checkoutOptions := store_working_directory.CheckoutOptions{
		CheckoutMode: c.CheckoutMode,
		Format:       zettel.Text{},
	}

	var checkoutResults zettel_checked_out.Set

	query := zettel_named.FilterIdSet{
		Set: ids,
		Or:  c.Or,
	}

	if checkoutResults, err = u.StoreWorkingDirectory().Checkout(checkoutOptions, query); err != nil {
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

	if _, err = openVimOp.Run(checkoutResults.ToSliceFilesZettelen()...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var readResults []zettel_checked_out.Zettel

	readOp := user_ops.ReadCheckedOut{
		Umwelt: u,
		OptionsReadExternal: store_working_directory.OptionsReadExternal{
			Format: zettel.Text{},
		},
	}

	fs := checkoutResults.ToSliceFilesZettelen()

	var possible store_working_directory.CwdFiles

	if possible, err = store_working_directory.MakeCwdFiles(u.Standort().Cwd(), fs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if readResults, err = readOp.RunMany(possible); err != nil {
		err = errors.Wrap(err)
		return
	}

	zettels := make([]zettel_external.Zettel, 0, len(readResults))

	for _, z := range readResults {
		zettels = append(zettels, z.External)
	}

	checkinOp := user_ops.Checkin{
		Umwelt: u,
		OptionsReadExternal: store_working_directory.OptionsReadExternal{
			Format: zettel.Text{},
		},
	}

	if _, err = checkinOp.Run(zettels...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
