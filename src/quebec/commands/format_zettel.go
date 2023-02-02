package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/lima/zettel_checked_out"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type FormatZettel struct {
	Format   string
	UTIGroup string
	Mode     zettel_checked_out.Mode
}

func init() {
	registerCommand(
		"format-zettel",
		func(f *flag.FlagSet) Command {
			c := &FormatZettel{
				Mode: zettel_checked_out.ModeZettelAndAkte,
			}

			f.Var(&c.Mode, "mode", "zettel, akte, or both")
			f.StringVar(&c.UTIGroup, "uti-group", "", "lookup format from UTI group")

			return c
		},
	)
}

func (c *FormatZettel) Run(u *umwelt.Umwelt, args ...string) (err error) {
	formatId := "text"
	var h kennung.Hinweis

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

	typKonfig := u.Konfig().GetApproximatedTyp(zt.Objekte.Typ).ApproximatedOrActual()

	var akteFormatter script_config.ScriptConfig

	if typKonfig != nil {
		actualFormatId := formatId
		var f script_config.WithOutputFormat
		ok := false

		if c.UTIGroup != "" {
			if g, ok := typKonfig.Objekte.Akte.FormatterUTIGroups[c.UTIGroup]; ok {
				if ft, ok := g.Map()[formatId]; ok {
					actualFormatId = ft
				}
			}
		}

		f, ok = typKonfig.Objekte.Akte.Formatters[actualFormatId]

		if ok {
			akteFormatter = f.ScriptConfig
		} else {
			err = errors.Normalf(
				"format '%s' for Typ '%s' not found",
				actualFormatId,
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
