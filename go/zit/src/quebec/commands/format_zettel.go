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
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type FormatZettel struct {
	Format string
	ids.RepoId
	UTIGroup string
	Mode     checkout_mode.Mode
}

func init() {
	registerCommand(
		"format-zettel",
		func(f *flag.FlagSet) Command {
			c := &FormatZettel{
				Mode: checkout_mode.MetadataAndBlob,
			}

			f.Var(&c.Mode, "mode", "metadata, blob, or both")

			f.Var(&c.RepoId, "kasten", "none or Browser")

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

func (c *FormatZettel) Run(u *env.Env, args ...string) (err error) {
	formatId := "text"

	var objectIdString string

	switch len(args) {
	case 1:
		objectIdString = args[0]

	case 2:
		formatId = args[0]
		objectIdString = args[1]

	default:
		err = errors.Errorf(
			"expected one or two input arguments, but got %d",
			len(args),
		)
		return
	}

	var zt *sku.Transacted

	if zt, err = c.getSku(u, objectIdString); err != nil {
		err = errors.Wrap(err)
		return
	}

	var blobFormatter script_config.RemoteScript

	if blobFormatter, err = c.getBlobFormatter(u, zt, formatId); err != nil {
		err = errors.Wrap(err)
		return
	}

	f := blob_store.MakeTextFormatterWithBlobFormatter(
		checkout_options.TextFormatterOptions{
			DoNotWriteEmptyDescription: true,
		},
		u.GetFSHome(),
		u.GetConfig(),
		blobFormatter,
	)

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
	u *env.Env,
	objectIdString string,
) (sk *sku.Transacted, err error) {
	b := u.MakeQueryBuilder(ids.MakeGenre(genres.Zettel))

	var qg *query.Group

	if qg, err = b.BuildQueryGroup(objectIdString); err != nil {
		err = errors.Wrap(err)
		return
	}

	var k *ids.ObjectId
	var s ids.Sigil

	if k, s, err = qg.GetExactlyOneObjectId(
		genres.Zettel,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if sk, err = u.GetStore().ReadTransactedFromObjectIdRepoIdSigil(
		k,
		c.RepoId,
		s,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *FormatZettel) getBlobFormatter(
	u *env.Env,
	zt *sku.Transacted,
	formatId string,
) (blobFormatter script_config.RemoteScript, err error) {
	if zt.GetType().IsEmpty() {
		ui.Log().Print("empty typ")
		return
	}

	var typKonfig *sku.Transacted

	if typKonfig, err = u.GetStore().ReadTransactedFromObjectId(zt.GetType()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var typeBlob *type_blobs.V0

	if typeBlob, err = u.GetStore().GetBlobStore().GetTypeV0().GetBlob(
		typKonfig.GetBlobSha(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	actualFormatId := formatId
	ok := false

	if c.UTIGroup == "" {
		blobFormatter, ok = typeBlob.Formatters[actualFormatId]

		if !ok {
			ui.Log().Print("no matching format id")
			blobFormatter = nil
			// TODO-P2 allow option to error on missing format
			// err = errors.Normalf("no format id %q", actualFormatId)
			// return
		}
	} else {
		var g type_blobs.FormatterUTIGroup
		g, ok = typeBlob.FormatterUTIGroups[c.UTIGroup]

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

		blobFormatter, ok = typeBlob.Formatters[actualFormatId]

		if !ok {
			ui.Log().Print("no matching format id")
			blobFormatter = nil
			// TODO-P2 allow option to error on missing format
			// err = errors.Normalf("no format id %q", actualFormatId)
			// return
		}
	}

	return
}
