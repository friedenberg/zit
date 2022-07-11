package commands

import (
	"flag"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/vim_cli_options_builder"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/juliett/user_ops"
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

func (c Edit) Run(u _Umwelt, args ...string) (err error) {
	checkoutOp := user_ops.Checkout{
		Umwelt: u,
		Options: _ZettelsCheckinOptions{
			IncludeAkte: c.IncludeAkte,
			Format:      _ZettelFormatsText{},
		},
	}

	var checkoutResults user_ops.CheckoutResults

	if checkoutResults, err = checkoutOp.Run(args...); err != nil {
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
		Umwelt:  u,
		Options: checkoutOp.Options,
	}

	var readResults user_ops.ReadCheckedOutResults

	readOp := user_ops.ReadCheckedOut{
		Umwelt:  u,
		Options: checkoutOp.Options,
	}

	if readResults, err = readOp.Run(checkoutResults.FilesZettelen...); err != nil {
		err = errors.Error(err)
		return
	}

	zettels := make([]stored_zettel.External, 0, len(readResults.Zettelen))

	for _, z := range readResults.Zettelen {
		zettels = append(zettels, z.External)
	}

	if _, err = checkinOp.Run(zettels...); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
