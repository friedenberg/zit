package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type FormatZettel struct {
	Format   string
	UTIGroup string
	Mode     sku.CheckoutMode
}

func init() {
	registerCommand(
		"format-zettel",
		func(f *flag.FlagSet) Command {
			c := &FormatZettel{
				Mode: sku.CheckoutModeObjekteAndAkte,
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

	var cwdFiles cwd.CwdFiles

	if cwdFiles, err = cwd.MakeCwdFilesAll(
		u.Konfig(),
		u.Standort().Cwd(),
		u.StoreObjekten(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	zt := &zettel.Transacted{}

	if e, ok := cwdFiles.GetZettel(h); ok {
		var ze zettel.External

		if ze, err = u.StoreObjekten().Zettel().ReadOneExternal(
			e,
			zt,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		zt = &zettel.Transacted{
			Objekte: ze.Objekte,
		}

		zt.Sku.Kennung = ze.Sku.Kennung
		zt.Sku.ObjekteSha = ze.Sku.ObjekteSha
		zt.Sku.AkteSha = ze.Sku.AkteSha
	} else {
		if zt, err = u.StoreObjekten().Zettel().ReadOne(&h); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	typKonfig := u.Konfig().GetApproximatedTyp(
		zt.GetTyp(),
	).ApproximatedOrActual()

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
				zt.GetTyp(),
			)

			return
		}
	}

	var format metadatei.TextFormatter

	if c.Mode.IncludesObjekte() {
		format = metadatei.MakeTextFormatterMetadateiInlineAkte(
			u.StoreObjekten(),
			akteFormatter,
		)
	} else {
		format = metadatei.MakeTextFormatterExcludeMetadatei(
			u.StoreObjekten(),
			akteFormatter,
		)
	}

	if err = u.Konfig().ApplyToMetadatei(zt); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO use cat or just write to stdout if no script instead of erroring
	if _, err = format.FormatMetadatei(u.Out(), zt); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
