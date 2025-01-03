package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type FormatObject struct {
	CheckoutMode checkout_mode.Mode // add test that says this is unused for stdin
	Stdin        bool               // switch to using `-`
	ids.RepoId
	UTIGroup string
}

func init() {
	registerCommand(
		"format-object",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &FormatObject{
				CheckoutMode: checkout_mode.BlobOnly,
			}

			f.BoolVar(&c.Stdin, "stdin", false, "Read object from stdin and use a Type directly")

			f.Var(&c.RepoId, "kasten", "none or Browser")

			f.StringVar(&c.UTIGroup, "uti-group", "", "lookup format from UTI group")

			f.Var(&c.CheckoutMode, "mode", "mode for checking out the zettel")

			return c
		},
	)
}

func (c *FormatObject) RunWithRepo(u *repo_local.Repo, args ...string) {
	if c.Stdin {
		if err := c.FormatFromStdin(u, args...); err != nil {
			u.CancelWithError(err)
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
		u.CancelWithErrorf(
			"expected one or two input arguments, but got %d",
			len(args),
		)
	}

	var object *sku.Transacted

	{
		var err error

		if object, err = u.GetSkuFromObjectId(objectIdString); err != nil {
			u.CancelWithError(err)
		}
	}

	tipe := object.GetType()

	{
		var err error

		if blobFormatter, err = u.GetBlobFormatter(
			tipe,
			formatId,
			c.UTIGroup,
		); err != nil {
			u.CancelWithError(err)
		}
	}

	f := blob_store.MakeTextFormatterWithBlobFormatter(
		u.GetRepoLayout(),
		checkout_options.TextFormatterOptions{
			DoNotWriteEmptyDescription: true,
		},
		u.GetConfig(),
		blobFormatter,
	)

	if err := u.GetStore().TryFormatHook(object); err != nil {
		u.CancelWithError(err)
	}

	if _, err := f.WriteStringFormatWithMode(
		u.GetUIFile(),
		object,
		c.CheckoutMode,
	); err != nil {
		u.CancelWithError(err)
	}
}

func (c *FormatObject) FormatFromStdin(
	u *repo_local.Repo,
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

	if blobFormatter, err = u.GetBlobFormatter(
		tipe,
		formatId,
		c.UTIGroup,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var wt io.WriterTo

	if wt, err = script_config.MakeWriterToWithStdin(
		blobFormatter,
		u.GetRepoLayout().MakeCommonEnv(),
		u.GetInFile(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = wt.WriteTo(u.GetUIFile()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
