package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("format-blob", &FormatBlob{})
}

type FormatBlob struct {
	command_components.LocalWorkingCopy

	complete command_components.Complete

	Stdin    bool
	UTIGroup string
}

func (cmd *FormatBlob) SetFlagSet(f *flag.FlagSet) {
	f.BoolVar(&cmd.Stdin, "stdin", false, "Read object from stdin and use a Type directly")

	f.StringVar(
		&cmd.UTIGroup,
		"uti-group",
		"",
		"lookup format from UTI group",
	)
}

func (cmd *FormatBlob) Complete(
	req command.Request,
	envLocal env_local.Env,
	commandLine command.CommandLine,
) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(req)

	args := commandLine.FlagsOrArgs[1:]

	if commandLine.InProgress != "" {
		args = args[:len(args)-1]
	}

	cmd.complete.CompleteObjects(
		req,
		localWorkingCopy,
		query.BuilderOptionDefaultGenres(genres.Zettel),
		args...,
	)
}

func (cmd *FormatBlob) Run(dep command.Request) {
	args := dep.PopArgs()
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

		if object, err = localWorkingCopy.GetZettelFromObjectId(objectIdString); err != nil {
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
			localWorkingCopy.CancelWithErrorAndFormat(
				err,
				"objectIdString: %q, Object: %q",
				objectIdString, sku.String(object),
			)
		}
	}

	format := typed_blob_store.MakeTextFormatterWithBlobFormatter(
		localWorkingCopy.GetEnvRepo(),
		checkout_options.TextFormatterOptions{
			DoNotWriteEmptyDescription: true,
		},
		localWorkingCopy.GetConfig(),
		blobFormatter,
		checkout_mode.None,
	)

	if err := localWorkingCopy.GetStore().TryFormatHook(object); err != nil {
		localWorkingCopy.CancelWithError(err)
	}

	if _, err := format.WriteStringFormatWithMode(
		localWorkingCopy.GetUIFile(),
		object,
		checkout_mode.BlobOnly,
	); err != nil {
		localWorkingCopy.CancelWithError(err)
	}
}

func (c *FormatBlob) FormatFromStdin(
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
		u.CancelWithError(err)
	}

	var wt io.WriterTo

	if wt, err = script_config.MakeWriterToWithStdin(
		blobFormatter,
		u.GetEnvRepo().MakeCommonEnv(),
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
