package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/hotel/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"golang.org/x/exp/maps"
)

type FormatBlob struct {
	Stdin  bool
	Format string
	ids.RepoId
	UTIGroup string
}

func init() {
	registerCommand(
		"format-blob",
		func(f *flag.FlagSet) Command {
			c := &FormatBlob{}

			f.BoolVar(&c.Stdin, "stdin", false, "Read object from stdin and use a Type directly")

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

func (c *FormatBlob) Run(u *env.Local, args ...string) (err error) {
	if c.Stdin {
		if err = c.FormatFromStdin(u, args...); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	var formatId string

	var objectIdString string
	var blobFormatter script_config.RemoteScript

	switch len(args) {
	case 2:
		formatId = args[1]
		fallthrough

	case 1:
		objectIdString = args[0]

	default:
		err = errors.Errorf(
			"expected one or two input arguments, but got %d",
			len(args),
		)
		return
	}

	var object *sku.Transacted

	if object, err = c.getSku(u, objectIdString); err != nil {
		err = errors.Wrap(err)
		return
	}

	tipe := object.GetType()

	if blobFormatter, err = c.getBlobFormatter(
		u,
		tipe,
		formatId,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	f := blob_store.MakeTextFormatterWithBlobFormatter(
		u.GetDirectoryLayout(),
		checkout_options.TextFormatterOptions{
			DoNotWriteEmptyDescription: true,
		},
		u.GetConfig(),
		blobFormatter,
	)

	if err = u.GetStore().TryFormatHook(object); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = f.WriteStringFormatWithMode(
		u.Out(),
		object,
		checkout_mode.BlobOnly,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *FormatBlob) FormatFromStdin(
	u *env.Local,
	args ...string,
) (err error) {
	formatId := "text"

	var blobFormatter script_config.RemoteScript
	var tipe ids.Type

	switch len(args) {
	case 1:
		if err = tipe.Set(args[0]); err != nil {
			err = errors.Wrap(err)
			return
		}

	case 2:
		formatId = args[0]
		if err = tipe.Set(args[1]); err != nil {
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

	if blobFormatter, err = c.getBlobFormatter(
		u,
		tipe,
		formatId,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var wt io.WriterTo

	if wt, err = script_config.MakeWriterToWithStdin(
		blobFormatter,
		u.GetDirectoryLayout().MakeCommonEnv(),
		u.In(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = wt.WriteTo(u.Out()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *FormatBlob) getSku(
	u *env.Local,
	objectIdString string,
) (sk *sku.Transacted, err error) {
	b := u.MakeQueryBuilder(ids.MakeGenre(genres.Zettel))

	var qg *query.Group

	if qg, err = b.BuildQueryGroup(objectIdString); err != nil {
		err = errors.Wrap(err)
		return
	}

	var e query.Executor

	if e, err = u.GetStore().MakeQueryExecutor(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	if sk, err = e.ExecuteExactlyOne(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *FormatBlob) getBlobFormatter(
	u *env.Local,
	tipe ids.Type,
	formatId string,
) (blobFormatter script_config.RemoteScript, err error) {
	if tipe.GetType().IsEmpty() {
		ui.Err().Print("empty type")
		return
	}

	var typeObject *sku.Transacted

	if typeObject, err = u.GetStore().ReadTransactedFromObjectId(
		tipe.GetType(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var typeBlob type_blobs.Blob

	if typeBlob, _, err = u.GetStore().GetBlobStore().GetType().ParseTypedBlob(
		typeObject.GetType(),
		typeObject.GetBlobSha(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer u.GetStore().GetBlobStore().GetType().PutTypedBlob(
		typeObject.GetType(),
		typeBlob,
	)

	ok := false

	if c.UTIGroup == "" {
		getBlobFormatter := func(formatId string) script_config.RemoteScript {
			var formatIds []string

			if formatId == "" {
				formatIds = []string{"text-edit", "text"}
			} else {
				formatIds = []string{formatId}
			}

			for _, formatId := range formatIds {
				blobFormatter, ok = typeBlob.GetFormatters()[formatId]

				if ok {
					return blobFormatter
				}
			}

			return nil
		}

		blobFormatter = getBlobFormatter(formatId)

		return
	}

	var g type_blobs.UTIGroup
	g, ok = typeBlob.GetFormatterUTIGroups()[c.UTIGroup]

	if !ok {
		err = errors.BadRequestf(
			"no uti group: %q. Available groups: %s",
			c.UTIGroup,
			maps.Keys(typeBlob.GetFormatterUTIGroups()),
		)
		return
	}

	ft, ok := g.Map()[formatId]

	if !ok {
		err = errors.Errorf(
			"no format id %q for uti group %q. Available groups: %q",
			formatId,
			c.UTIGroup,
			maps.Keys(g.Map()),
		)

		return
	}

	formatId = ft

	blobFormatter, ok = typeBlob.GetFormatters()[formatId]

	if !ok {
		ui.Err().Print("no matching format id")
		blobFormatter = nil
		// TODO-P2 allow option to error on missing format
		// err = errors.Normalf("no format id %q", actualFormatId)
		// return

		return
	}

	return
}
