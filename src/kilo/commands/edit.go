package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	checkout_store "github.com/friedenberg/zit/src/hotel/store_checkout"
	"github.com/friedenberg/zit/src/india/store_with_lock"
	"github.com/friedenberg/zit/src/juliett/user_ops"
)

type Edit struct {
	IncludeAkte bool
}

func init() {
	registerCommand(
		"edit",
		func(f *flag.FlagSet) Command {
			c := &Edit{}

			f.BoolVar(&c.IncludeAkte, "include-akte", true, "check out and open the akte")

			return c
		},
	)
}

func (c Edit) Run(u *umwelt.Umwelt, args ...string) (err error) {
	checkoutOp := user_ops.Checkout{
		Umwelt: u,
		Options: checkout_store.CheckinOptions{
			IncludeAkte: c.IncludeAkte,
			Format:      zettel_formats.Text{},
		},
	}

	var hins []hinweis.Hinweis

	if hins, err = (user_ops.GetHinweisenFromArgs{}).RunMany(args...); err != nil {
		err = errors.Error(err)
		return
	}

	var checkoutResults user_ops.CheckoutResults

	var s store_with_lock.Store

	if s, err = store_with_lock.New(u); err != nil {
		err = errors.Error(err)
		return
	}

	if checkoutResults, err = checkoutOp.RunManyHinweisen(s, hins...); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.Flush(); err != nil {
		err = errors.Error(err)
		return
	}

	if err = (user_ops.OpenFiles{}).Run(checkoutResults.FilesAkten...); err != nil {
		err = errors.Error(err)
		return
	}

	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithCursorLocation(2, 3).
			WithFileType("zit-zettel").
			WithInsertMode().
			Build(),
	}

	if _, err = openVimOp.Run(checkoutResults.FilesZettelen...); err != nil {
		err = errors.Error(err)
		return
	}

	checkinOp := user_ops.Checkin{
		Umwelt:  s.Umwelt,
		Options: checkoutOp.Options,
	}

	var readResults user_ops.ReadCheckedOutResults

	readOp := user_ops.ReadCheckedOut{
		Umwelt:  s.Umwelt,
		Options: checkoutOp.Options,
	}

	if s, err = store_with_lock.New(u); err != nil {
		err = errors.Error(err)
		return
	}

	if readResults, err = readOp.RunManyStrings(s, checkoutResults.FilesZettelen...); err != nil {
		err = errors.Error(err)
		return
	}

	zettels := make([]stored_zettel.External, 0, len(readResults.Zettelen))

	for _, z := range readResults.Zettelen {
		zettels = append(zettels, z.External)
	}

	if _, err = checkinOp.Run(s, zettels...); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.Flush(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
