package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/mike/umwelt"
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

	if f, err = files.Open(args[0]); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer files.Close(f)

	format := zettel.Text{}

	var cz zettel_checked_out.Zettel

	if cz, err = u.StoreWorkingDirectory().Read(args[0]); err != nil {
		err = errors.Wrap(err)
		return
	}

	var formatter konfig.RemoteScript

	if typKonfig, ok := u.Konfig().Typen[cz.External.Named.Stored.Zettel.Typ.String()]; ok {
		formatter = typKonfig.FormatScript
	}

	if err = cz.External.Named.Stored.Zettel.ApplyKonfig(u.Konfig()); err != nil {
		err = errors.Wrap(err)
		return
	}

	ctx := zettel.FormatContextWrite{
		Zettel:            cz.External.Named.Stored.Zettel,
		IncludeAkte:       true,
		AkteReaderFactory: u.StoreObjekten(),
		FormatScript:      formatter,
		Out:               os.Stdout,
	}

	if _, err = format.WriteTo(ctx); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
