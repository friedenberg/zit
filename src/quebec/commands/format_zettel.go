package commands

import (
	"flag"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/zettel_checked_out"
	"github.com/friedenberg/zit/src/oscar/umwelt"
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

	defer errors.Deferred(&err, f.Close)

	//TODO-P4 read the format type (akte included or not) and determine
	//appropriately
	format := zettel.MakeObjekteTextFormatterIncludeAkte(
		u.Konfig(),
		u.StoreObjekten(),
		nil,
	)

	var cz zettel_checked_out.Zettel

	if cz, err = u.StoreWorkingDirectory().Read(args[0]); err != nil {
		err = errors.Wrap(err)
		return
	}

	typKonfig := u.Konfig().GetTyp(cz.External.Objekte.Typ)

	if typKonfig != nil {
		if f, ok := typKonfig.Objekte.Akte.Actions["format"]; ok {
			format.AkteFormatter = f.ScriptConfig
		}
	}

	if err = cz.External.Objekte.ApplyKonfig(u.Konfig()); err != nil {
		err = errors.Wrap(err)
		return
	}

	ctx := zettel.ObjekteFormatterContext{
		Zettel:      cz.External.Objekte,
		IncludeAkte: u.Konfig().IsInlineTyp(cz.External.Objekte.Typ),
	}

	if _, err = format.Format(u.Out(), &ctx); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
