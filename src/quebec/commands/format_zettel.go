package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/charlie/script_config"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/india/external"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/lima/cwd"
	"github.com/friedenberg/zit/src/lima/store_objekten"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type FormatZettel struct {
	Format   string
	UTIGroup string
	Mode     checkout_mode.Mode
}

func init() {
	registerCommand(
		"format-zettel",
		func(f *flag.FlagSet) Command {
			c := &FormatZettel{
				Mode: checkout_mode.ModeObjekteAndAkte,
			}

			f.Var(&c.Mode, "mode", "zettel, akte, or both")
			f.StringVar(
				&c.UTIGroup,
				"uti-group",
				"",
				"lookup format from UTI group",
			)

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
		err = errors.Errorf(
			"expected one or two input arguments, but got %d",
			len(args),
		)
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

	var zt *transacted.Zettel

	if zt, err = u.StoreObjekten().Zettel().ReadOne(&h); err != nil {
		err = errors.Wrap(err)
		return
	}

	if e, ok := cwdFiles.GetZettel(&h); ok {
		var ze external.Zettel

		ze, err = u.StoreObjekten().Zettel().ReadOneExternal(e, zt)

		switch {
		case store_objekten.IsErrExternalAkteExtensionMismatch(err):
			err = nil

		case err != nil:
			err = errors.Wrap(err)
			return

		default:
			// TODO-P1 switch to methods on Transacted and External
			zt.SetMetadatei(ze.GetMetadatei())
			zt.Kennung = ze.GetKennung()
			zt.SetObjekteSha(ze.GetObjekteSha())
			zt.SetAkteSha(ze.GetAkteSha())
		}
	}

	typKonfig := u.Konfig().GetApproximatedTyp(
		zt.GetTyp(),
	).ApproximatedOrActual()

	var akteFormatter script_config.RemoteScript

	if typKonfig == nil {
		panic("typ konfig was nil")
	}

	var typAkte *typ_akte.V0

	if typAkte, err = u.StoreObjekten().Typ().GetAkte(typKonfig.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	actualFormatId := formatId
	ok := false

	if c.UTIGroup != "" {
		g, ok := typAkte.FormatterUTIGroups[c.UTIGroup]

		if !ok {
			err = errors.Errorf("no uti group: %q", c.UTIGroup)
			return
		}

		ft, ok := g.Map()[formatId]

		if !ok {
			err = errors.Errorf(
				"no format id %q for uti group %q",
				formatId,
				c.UTIGroup,
			)

			return
		}

		actualFormatId = ft
	}

	akteFormatter, ok = typAkte.Formatters[actualFormatId]

	if !ok {
		err = errors.Errorf("no format id %q", actualFormatId)
		return
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

	if err = u.Konfig().ApplyToNewMetadatei(zt, u.StoreObjekten().Typ()); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO-P1 use cat or just write to stdout if no script instead of erroring
	if _, err = format.FormatMetadatei(u.Out(), zt); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
