package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type FormatZettel struct {
	Format   string
	Kasten   ids.RepoId
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

			f.Var(&c.Kasten, "kasten", "none or Chrome")

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

	var zt *sku.Transacted

	if zt, err = c.getSku(u, kennungString); err != nil {
		err = errors.Wrap(err)
		return
	}

	var akteFormatter script_config.RemoteScript

	if akteFormatter, err = c.getAkteFormatter(u, zt, formatId); err != nil {
		err = errors.Wrap(err)
		return
	}

	f := blob_store.MakeTextFormatterWithBlobFormatter(
		checkout_options.TextFormatterOptions{
			DoNotWriteEmptyBezeichnung: true,
		},
		u.Standort(),
		u.GetKonfig(),
		akteFormatter,
	)

	if err = u.GetKonfig().ApplyToNewMetadata(
		zt,
		u.GetStore().GetAkten().GetTypeV0(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetStore().TryFormatHook(zt); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = f.WriteStringFormatWithMode(u.Out(), zt, c.Mode); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *FormatZettel) getSku(
	u *umwelt.Umwelt,
	kennungString string,
) (sk *sku.Transacted, err error) {
	b := u.MakeQueryBuilder(ids.MakeGenre(genres.Zettel))

	var qg *query.Group

	if qg, err = b.BuildQueryGroup(kennungString); err != nil {
		err = errors.Wrap(err)
		return
	}

	var k *ids.ObjectId
	var s ids.Sigil

	if k, s, err = qg.GetExactlyOneKennung(
		genres.Zettel,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if sk, err = u.GetStore().ReadTransactedFromKennungKastenSigil(
		k,
		c.Kasten,
		s,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *FormatZettel) getAkteFormatter(
	u *umwelt.Umwelt,
	zt *sku.Transacted,
	formatId string,
) (akteFormatter script_config.RemoteScript, err error) {
	if zt.GetType().IsEmpty() {
		ui.Log().Print("empty typ")
		return
	}

	var typKonfig *sku.Transacted

	if typKonfig, err = u.GetStore().ReadTransactedFromObjectId(zt.GetType()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var typAkte *type_blobs.V0

	if typAkte, err = u.GetStore().GetAkten().GetTypeV0().GetBlob(
		typKonfig.GetBlobSha(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	actualFormatId := formatId
	ok := false

	if c.UTIGroup == "" {
		akteFormatter, ok = typAkte.Formatters[actualFormatId]

		if !ok {
			ui.Log().Print("no matching format id")
			akteFormatter = nil
			// TODO-P2 allow option to error on missing format
			// err = errors.Normalf("no format id %q", actualFormatId)
			// return
		}
	} else {
		var g type_blobs.FormatterUTIGroup
		g, ok = typAkte.FormatterUTIGroups[c.UTIGroup]

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

		akteFormatter, ok = typAkte.Formatters[actualFormatId]

		if !ok {
			ui.Log().Print("no matching format id")
			akteFormatter = nil
			// TODO-P2 allow option to error on missing format
			// err = errors.Normalf("no format id %q", actualFormatId)
			// return
		}
	}

	return
}
