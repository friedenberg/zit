package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_external"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/india/store_working_directory"
	"github.com/friedenberg/zit/src/juliett/umwelt"
	"github.com/friedenberg/zit/src/kilo/user_ops"
)

type Edit struct {
	store_working_directory.CheckoutMode
}

func init() {
	registerCommand(
		"edit",
		func(f *flag.FlagSet) Command {
			c := &Edit{
				CheckoutMode: store_working_directory.CheckoutModeZettelOnly,
			}

			f.Var(&c.CheckoutMode, "mode", "mode for checking out the zettel")

			return c
		},
	)
}

func (c Edit) Run(u *umwelt.Umwelt, args ...string) (err error) {
	checkoutOptions := store_working_directory.CheckoutOptions{
		CheckoutMode: c.CheckoutMode,
		Format:       zettel.Text{},
	}

	checkoutOp := user_ops.Checkout{
		Umwelt:          u,
		CheckoutOptions: checkoutOptions,
	}

	var hins []hinweis.Hinweis

	op := user_ops.GetHinweisenFromArgs{
		Umwelt: u,
	}

	if hins, err = op.RunMany(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	var checkoutResults zettel_checked_out.Set

	if checkoutResults, err = checkoutOp.RunManyHinweisen(hins...); err != nil {
		err = errors.Wrap(err)
		return
	}

	switch c.CheckoutMode {
	case store_working_directory.CheckoutModeAkteOnly:
	case store_working_directory.CheckoutModeZettelAndAkte:
		if err = (user_ops.OpenFiles{}).Run(checkoutResults.ToSliceFilesAkten()...); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
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

	var possible store_working_directory.CwdFiles

	for _, f := range checkoutResults.ToSliceFilesZettelen() {
		possible.Zettelen = append(possible.Zettelen, f)
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
