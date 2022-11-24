package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/echo/konfig"
	"github.com/friedenberg/zit/src/hotel/zettel"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/november/umwelt"
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

	typKonfig := u.Konfig().GetTyp(cz.External.Named.Stored.Objekte.Typ.String())

	if typKonfig != nil {
		if f, ok := typKonfig.Actions["format"]; ok {
			formatter = f
		}
	}

	if err = cz.External.Named.Stored.Objekte.ApplyKonfig(u.Konfig()); err != nil {
		err = errors.Wrap(err)
		return
	}

	ctx := zettel.FormatContextWrite{
		Zettel:            cz.External.Named.Stored.Objekte,
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
