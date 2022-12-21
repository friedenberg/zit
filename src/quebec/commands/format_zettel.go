package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/zettel_checked_out"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type FormatZettel struct {
	Format string
	Mode   zettel_checked_out.Mode
}

func init() {
	registerCommand(
		"format-zettel",
		func(f *flag.FlagSet) Command {
			c := &FormatZettel{
				Mode: zettel_checked_out.ModeZettelAndAkte,
			}

			f.Var(&c.Mode, "mode", "zettel, akte, or both")

			return c
		},
	)
}

// TODO-P1 migrate to using store_fs
func (c *FormatZettel) Run(u *umwelt.Umwelt, args ...string) (err error) {
	formatId := "text"
	var h hinweis.Hinweis

	switch len(args) {
	case 1:
		if err = h.Set(args[0]); err != nil {
			err = errors.Wrap(err)
			return
		}

	case 2:
		formatId = args[0]

		if err = h.Set(args[1]); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = errors.Errorf("expected one or two input arguments, but got %d", len(args))
		return
	}

	var zt *zettel.Transacted

	if zt, err = u.StoreWorkingDirectory().ReadOne(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	typKonfig := u.Konfig().GetTyp(zt.Objekte.Typ)

	var akteFormatter script_config.ScriptConfig

	if typKonfig != nil {
		if f, ok := typKonfig.Objekte.Akte.Formatters[formatId]; ok {
			akteFormatter = f.ScriptConfig
		} else {
			err = errors.Normalf(
				"format '%s' for Typ '%s' not found",
				formatId,
				zt.Objekte.Typ,
			)

			return
		}
	}

	var format zettel.ObjekteFormatter

	if c.Mode.IncludesZettel() {
		format = zettel.MakeObjekteTextFormatterIncludeAkte(
			u.Standort(),
			u.Konfig(),
			u.StoreObjekten(),
			akteFormatter,
		)
	} else {
		format = zettel.MakeObjekteTextFormatterExcludeMetadatei(
			u.Standort(),
			u.Konfig(),
			u.StoreObjekten(),
			akteFormatter,
		)
	}

	//TODO-P1 what is this?
	if err = zt.Objekte.ApplyKonfig(u.Konfig()); err != nil {
		err = errors.Wrap(err)
		return
	}

	ctx := zettel.ObjekteFormatterContext{
		Zettel:      zt.Objekte,
		IncludeAkte: u.Konfig().IsInlineTyp(zt.Objekte.Typ) && c.Mode.IncludesAkte(),
	}

	if _, err = format.Format(u.Out(), &ctx); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
