package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type FormatZettel struct {
}

func init() {
	registerCommand(
		"format-zettel",
		func(f *flag.FlagSet) Command {
			c := &FormatZettel{}

			return c
		},
	)
}

func (c *FormatZettel) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) != 1 {
		err = errors.Errorf("expected exactly one input argument")
		return
	}

	// stdoutIsTty := open_file_guard.IsTty(os.Stdout)
	// stdinIsTty := open_file_guard.IsTty(os.Stdin)

	// if !stdinIsTty && !stdoutIsTty {
	// 	logz.Print("neither stdin or stdout is a tty")
	// 	logz.Print("generate organize, read from stdin, commit")

	var f *os.File

	if f, err = open_file_guard.Open(args[0]); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(f)

	format := zettel_formats.Text{}

	// checkinOptions := zettels.CheckinOptions{
	// 	IgnoreMissingHinweis: true,
	// 	IncludeAkte:          true,
	// 	Format:               format,
	// }

	// readOp := user_ops.ReadCheckedOut{
	// 	Umwelt:  u,
	// 	Options: checkinOptions,
	// }

	// var z stored_zettel.CheckedOut

	// if z, err = readOp.RunOne(args[0]); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	var store store_with_lock.Store

	if store, err = store_with_lock.New(u); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	var external stored_zettel.External

	if external, err = store.CheckoutStore().Read(args[0]); err != nil {
		err = errors.Error(err)
		return
	}

	ctx := zettel.FormatContextWrite{
		Zettel:            external.Zettel,
		IncludeAkte:       false,
		AkteReaderFactory: store.Zettels(),
		Out:               os.Stdout,
	}

	if _, err = format.WriteTo(ctx); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
