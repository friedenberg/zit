package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/delta/konfig"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/golf/zettel_formats"
	"github.com/friedenberg/zit/src/india/store_with_lock"
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

	var f *os.File

	if f, err = open_file_guard.Open(args[0]); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(f)

	format := zettel_formats.Text{}

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

	var formatter konfig.RemoteScript

	if typKonfig, ok := u.Konfig.Typen[external.Named.Stored.Zettel.Typ.String()]; ok {
		formatter = typKonfig.FormatScript
	}

	ctx := zettel.FormatContextWrite{
		Zettel:            external.Named.Stored.Zettel,
		IncludeAkte:       true,
		AkteReaderFactory: store.StoreObjekten(),
		FormatScript:      formatter,
		Out:               os.Stdout,
	}

	if _, err = format.WriteTo(ctx); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
