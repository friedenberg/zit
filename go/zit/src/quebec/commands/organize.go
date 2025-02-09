package commands

import (
	"flag"
	"os"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/go/zit/src/bravo/organize_text_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

func init() {
	command.Register(
		"organize",
		&Organize{
			Flags: organize_text.MakeFlags(),
		},
	)
}

// Refactor and fold components into userops
type Organize struct {
	command_components.LocalWorkingCopyWithQueryGroup

	organize_text.Flags
	Mode organize_text_mode.Mode

	Filter script_value.ScriptValue
}

func (c *Organize) SetFlagSet(f *flag.FlagSet) {
	c.LocalWorkingCopyWithQueryGroup.SetFlagSet(f)

	c.Flags.SetFlagSet(f)

	f.Var(
		&c.Filter,
		"filter",
		"a script to run for each file to transform it the standard zettel format",
	)

	f.Var(&c.Mode, "mode", "mode used for handling stdin and stdout")
}

func (c *Organize) ModifyBuilder(b *query.Builder) {
	b.
		WithRequireNonEmptyQuery()
}

func (c *Organize) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
		genres.Tag,
		genres.Type,
	)
}

func (cmd *Organize) Run(req command.Request) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(req)
	envWorkspace := localWorkingCopy.GetEnvWorkspace()

	queryGroup := cmd.MakeQueryGroup(
		req,
		query.BuilderOptionsOld(
			cmd,
			query.BuilderOptionWorkspace{Env: envWorkspace},
			query.BuilderOptionDefaultGenres(genres.Zettel),
			query.BuilderOptionDefaultSigil(ids.SigilLatest),
		),
		localWorkingCopy,
		req.Args(),
	)

	localWorkingCopy.ApplyToOrganizeOptions(&cmd.Options)

	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Repo: localWorkingCopy,
		Options: localWorkingCopy.MakeOrganizeOptionsWithQueryGroup(
			cmd.Flags,
			queryGroup,
		),
	}

	types := query.GetTypes(queryGroup)

	if types.Len() == 1 {
		createOrganizeFileOp.Type = types.Any()
	}

	skus := sku.MakeSkuTypeSetMutable()
	var l sync.Mutex

	if err := localWorkingCopy.GetStore().QueryTransactedAsSkuType(
		queryGroup,
		func(co sku.SkuType) (err error) {
			l.Lock()
			defer l.Unlock()

			return skus.Add(co.Clone())
		},
	); err != nil {
		localWorkingCopy.CancelWithError(err)
	}

	createOrganizeFileOp.Skus = skus

	if createOrganizeFileOp.Skus.Len() == 0 && createOrganizeFileOp.TagSet.Len() == 0 {
		createOrganizeFileOp.TagSet = ids.MakeTagSet(
			localWorkingCopy.GetEnvWorkspace().GetDefaults().GetTags()...,
		)
	}

	switch cmd.Mode {
	case organize_text_mode.ModeCommitDirectly:
		ui.Log().Print("neither stdin or stdout is a tty")
		ui.Log().Print("generate organize, read from stdin, commit")

		var createOrganizeFileResults *organize_text.Text

		var f *os.File

		{
			var err error

			if f, err = files.TempFileWithPattern(
				"*." + localWorkingCopy.GetConfig().GetFileExtensions().GetFileExtensionOrganize(),
			); err != nil {
				localWorkingCopy.CancelWithError(err)
			}
		}

		defer localWorkingCopy.MustClose(f)

		{
			var err error

			if createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(
				f,
			); err != nil {
				localWorkingCopy.CancelWithError(err)
			}
		}

		var organizeText *organize_text.Text

		readOrganizeTextOp := user_ops.ReadOrganizeFile{}

		{
			var err error

			if organizeText, err = readOrganizeTextOp.Run(
				localWorkingCopy,
				os.Stdin,
				organize_text.NewMetadata(queryGroup.RepoId),
			); err != nil {
				localWorkingCopy.CancelWithError(err)
			}
		}

		if _, err := localWorkingCopy.LockAndCommitOrganizeResults(
			organize_text.OrganizeResults{
				Before:     createOrganizeFileResults,
				After:      organizeText,
				Original:   skus,
				QueryGroup: queryGroup,
			},
		); err != nil {
			localWorkingCopy.CancelWithError(err)
		}

	case organize_text_mode.ModeOutputOnly:
		ui.Log().Print("generate organize file and write to stdout")
		if _, err := createOrganizeFileOp.RunAndWrite(os.Stdout); err != nil {
			localWorkingCopy.CancelWithError(err)
		}

	case organize_text_mode.ModeInteractive:
		ui.Log().Print(
			"generate temp file, write organize, open vim to edit, commit results",
		)
		var createOrganizeFileResults *organize_text.Text

		var f *os.File

		{
			var err error

			if f, err = localWorkingCopy.GetEnvRepo().GetTempLocal().FileTempWithTemplate(
				"*." + localWorkingCopy.GetConfig().GetFileExtensions().GetFileExtensionOrganize(),
			); err != nil {
				localWorkingCopy.CancelWithError(err)
			}

			defer localWorkingCopy.MustClose(f)
		}

		{
			var err error

			if createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(
				f,
			); err != nil {
				localWorkingCopy.CancelWithErrorAndFormat(err, "Organize File: %q", f.Name())
			}
		}

		var organizeText *organize_text.Text

		{
			var err error

			if organizeText, err = cmd.readFromVim(
				localWorkingCopy,
				f.Name(),
				createOrganizeFileResults,
				queryGroup,
			); err != nil {
				localWorkingCopy.CancelWithErrorAndFormat(err, "Organize File: %q", f.Name())
			}
		}

		if _, err := localWorkingCopy.LockAndCommitOrganizeResults(
			organize_text.OrganizeResults{
				Before:     createOrganizeFileResults,
				After:      organizeText,
				Original:   skus,
				QueryGroup: queryGroup,
			},
		); err != nil {
			localWorkingCopy.CancelWithError(err)
		}

	default:
		localWorkingCopy.CancelWithErrorf("unknown mode")
	}
}

func (c Organize) readFromVim(
	repo *local_working_copy.Repo,
	path string,
	results *organize_text.Text,
	queryGroup *query.Query,
) (ot *organize_text.Text, err error) {
	openVimOp := user_ops.OpenEditor{
		VimOptions: vim_cli_options_builder.New().
			WithFileType("zit-organize").
			Build(),
	}

	if err = openVimOp.Run(repo, path); err != nil {
		err = errors.Wrap(err)
		return
	}

	readOrganizeTextOp := user_ops.ReadOrganizeFile{}

	if ot, err = readOrganizeTextOp.RunWithPath(
		repo,
		path,
		queryGroup.RepoId,
	); err != nil {
		if c.handleReadChangesError(repo, err) {
			err = nil
			ot, err = c.readFromVim(repo, path, results, queryGroup)
		} else {
			ui.Err().Printf("aborting organize")
			return
		}
	}

	return
}

// TODO migrate to using errors.Retryable
func (cmd Organize) handleReadChangesError(
	envUI env_ui.Env,
	err error,
) (tryAgain bool) {
	var errorRead organize_text.ErrorRead

	if err != nil && !errors.As(err, &errorRead) {
		ui.Err().Printf("unrecoverable organize read failure: %s", err)
		tryAgain = false
		return
	}

	tryAgain = envUI.Retry("reading changes failed", "edit and retry?", err)

	return
}
