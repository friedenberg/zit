package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	registerCommand(
		"format-object",
		&FormatObject{
			CheckoutMode: checkout_mode.BlobOnly,
		},
	)
}

type FormatObject struct {
	command_components.LocalWorkingCopy

	CheckoutMode checkout_mode.Mode // add test that says this is unused for stdin
	Stdin        bool               // switch to using `-`
	ids.RepoId
	UTIGroup string
}

func (cmd *FormatObject) SetFlagSet(f *flag.FlagSet) {
	f.BoolVar(&cmd.Stdin, "stdin", false, "Read object from stdin and use a Type directly")

	f.Var(&cmd.RepoId, "kasten", "none or Browser")

	f.StringVar(&cmd.UTIGroup, "uti-group", "", "lookup format from UTI group")

	f.Var(&cmd.CheckoutMode, "mode", "mode for checking out the zettel")
}

func (cmd *FormatObject) Run(dep command.Dep) {
	args := dep.Args()
	localWorkingCopy := cmd.MakeLocalWorkingCopy(dep)

	if cmd.Stdin {
		if err := cmd.FormatFromStdin(localWorkingCopy, args...); err != nil {
			localWorkingCopy.CancelWithError(err)
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
		localWorkingCopy.CancelWithErrorf(
			"expected one or two input arguments, but got %d",
			len(args),
		)
	}

	var object *sku.Transacted

	{
		var err error

		if object, err = localWorkingCopy.GetSkuFromObjectId(objectIdString); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	}

	tipe := object.GetType()

	{
		var err error

		if blobFormatter, err = localWorkingCopy.GetBlobFormatter(
			tipe,
			formatId,
			cmd.UTIGroup,
		); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	}

	f := blob_store.MakeTextFormatterWithBlobFormatter(
		localWorkingCopy.GetRepoLayout(),
		checkout_options.TextFormatterOptions{
			DoNotWriteEmptyDescription: true,
		},
		localWorkingCopy.GetConfig(),
		blobFormatter,
	)

	if err := localWorkingCopy.GetStore().TryFormatHook(object); err != nil {
		localWorkingCopy.CancelWithError(err)
	}

	if _, err := f.WriteStringFormatWithMode(
		localWorkingCopy.GetUIFile(),
		object,
		cmd.CheckoutMode,
	); err != nil {
		var errBlobFormatterFailed *object_metadata.ErrBlobFormatterFailed

		if errors.As(err, &errBlobFormatterFailed) {
			localWorkingCopy.CancelWithError(errBlobFormatterFailed)
			// err = nil
			// ui.Err().Print(errExit)
		} else {
			localWorkingCopy.CancelWithError(err)
		}
	}
}

func (c *FormatObject) FormatFromStdin(
	u *local_working_copy.Repo,
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
