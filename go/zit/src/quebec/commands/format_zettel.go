package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/query"
	"code.linenisgreat.com/zit/src/juliett/objekte"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
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

	var kennungString string

	switch len(args) {
	case 1:
		kennungString = args[0]

	case 2:
		formatId = args[0]
		kennungString = args[1]

	default:
		err = errors.Errorf(
			"expected one or two input arguments, but got %d",
			len(args),
		)
		return
	}

	b := u.MakeMetaIdSetWithoutExcludedHidden(kennung.MakeGattung(gattung.Zettel))

	var qg *query.Group

	if qg, err = b.BuildQueryGroup(kennungString); err != nil {
		err = errors.Wrap(err)
		return
	}

	var k *kennung.Kennung2
	var s kennung.Sigil

	if k, s, err = qg.GetExactlyOneKennung(
		gattung.Zettel,
		u.GetStore().GetCwdFiles(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var zt *sku.Transacted

	if zt, err = u.GetStore().ReadOneSigil(k, s); err != nil {
		err = errors.Wrap(err)
		return
	}

	typKonfig := u.Konfig().GetApproximatedTyp(
		zt.GetTyp(),
	).ApproximatedOrActual()

	var akteFormatter script_config.RemoteScript

	if typKonfig == nil {
		panic("typ konfig was nil")
	}

	var typAkte *typ_akte.V0

	if typAkte, err = u.GetStore().GetAkten().GetTypV0().GetAkte(
		typKonfig.GetAkteSha(),
	); err != nil {
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
		akteFormatter = nil
		// TODO-P2 allow option to error on missing format
		// err = errors.Normalf("no format id %q", actualFormatId)
		// return
	}

	f := objekte.MakeTextFormatterWithAkteFormatter(
		u.Standort(),
		u.Konfig(),
		akteFormatter,
	)

	if err = u.Konfig().ApplyToNewMetadatei(
		zt,
		u.GetStore().GetAkten().GetTypV0(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = f.WriteStringFormatWithMode(u.Out(), zt, c.Mode); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
