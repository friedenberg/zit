package user_ops

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type Organize2 struct {
	*local_working_copy.Repo
	organize_text.Metadata
}

func (op Organize2) Run(
	skus sku.CheckedOutMutableSet,
) (organizeResults organize_text.OrganizeResults, err error) {
	organizeResults.Original = skus

	organizeFlags := organize_text.MakeFlagsWithMetadata(op.Metadata)
	op.ApplyToOrganizeOptions(&organizeFlags.Options)
	organizeFlags.Skus = skus

	createOrganizeFileOp := CreateOrganizeFile{
		Repo: op.Repo,
		Options: op.Repo.MakeOrganizeOptionsWithOrganizeMetadata(
			organizeFlags,
			op.Metadata,
		),
	}

	var f *os.File

	if f, err = op.GetEnvRepo().GetTempLocal().FileTempWithTemplate(
		"*." + op.GetConfig().GetFileExtensions().GetFileExtensionOrganize(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if organizeResults.Before, err = createOrganizeFileOp.RunAndWrite(
		f,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO refactor into common vim processing loop
	for {
		openVimOp := OpenEditor{
			VimOptions: vim_cli_options_builder.New().
				WithFileType("zit-organize").
				Build(),
		}

		if err = openVimOp.Run(op.Repo, f.Name()); err != nil {
			err = errors.Wrap(err)
			return
		}

		// if err = op.Reset(); err != nil {
		// 	err = errors.Wrap(err)
		// 	return
		// }

		readOrganizeTextOp := ReadOrganizeFile{}

		if _, err = f.Seek(0, io.SeekStart); err != nil {
			err = errors.Wrap(err)
			return
		}

		if organizeResults.After, err = readOrganizeTextOp.Run(
			op.Repo,
			f,
			organize_text.NewMetadataWithOptionCommentLookup(
				organizeResults.Before.Metadata.RepoId,
				op.GetPrototypeOptionComments(),
			),
		); err != nil {
			if op.handleReadChangesError(op.Repo, err) {
				err = nil
				continue
			} else {
				ui.Err().Printf("aborting organize")
				return
			}
		}

		break
	}

	return
}

func (cmd Organize2) handleReadChangesError(
	envUI env_ui.Env,
	err error,
) (tryAgain bool) {
	var errorRead organize_text.ErrorRead

	if err != nil && !errors.As(err, &errorRead) {
		ui.Err().Printf("unrecoverable organize read failure: %s", err)
		tryAgain = false
		return
	}

	return envUI.Retry(
		"reading changes failed",
		"edit and try again?",
		err,
	)
}
