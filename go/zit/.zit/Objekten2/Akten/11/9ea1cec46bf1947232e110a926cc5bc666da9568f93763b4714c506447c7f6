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
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

// Refactor and fold components into userops
type Organize struct {
	organize_text.Flags
	Mode organize_text_mode.Mode

	Filter script_value.ScriptValue
}

func init() {
	registerCommandWithQuery(
		"organize",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Organize{
				Flags: organize_text.MakeFlags(),
			}

			f.Var(
				&c.Filter,
				"filter",
				"a script to run for each file to transform it the standard zettel format",
			)

			f.Var(&c.Mode, "mode", "mode used for handling stdin and stdout")

			c.AddToFlagSet(f)

			return c
		},
	)
}

func (c *Organize) ModifyBuilder(b *query.Builder) {
	b.
		WithDefaultSigil(ids.SigilLatest).
		WithDefaultGenres(ids.MakeGenre(genres.Zettel)).
		WithRequireNonEmptyQuery()
}

func (c *Organize) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
		genres.Tag,
		genres.Type,
	)
}

func (c *Organize) RunWithQuery(
	u *env.Local,
	qg *query.Group,
) (err error) {
	u.ApplyToOrganizeOptions(&c.Options)

	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Local: u,
		Options: u.MakeOrganizeOptionsWithQueryGroup(
			c.Flags,
			qg,
		),
	}

	typen := qg.GetTypes()

	if typen.Len() == 1 {
		createOrganizeFileOp.Type = typen.Any()
	}

	skus := sku.MakeSkuTypeSetMutable()
	var l sync.RWMutex

	if err = u.GetStore().QueryTransactedAsSkuType(
		qg,
		func(co sku.SkuType) (err error) {
			l.Lock()
			defer l.Unlock()

			return skus.Add(co.Clone())
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	createOrganizeFileOp.Skus = skus

	switch c.Mode {
	case organize_text_mode.ModeCommitDirectly:
		ui.Log().Print("neither stdin or stdout is a tty")
		ui.Log().Print("generate organize, read from stdin, commit")

		var createOrganizeFileResults *organize_text.Text

		var f *os.File

		if f, err = files.TempFileWithPattern(
			"*." + u.GetConfig().GetFileExtensions().GetFileExtensionOrganize(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, f)

		if createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(
			f,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		var organizeText *organize_text.Text

		readOrganizeTextOp := user_ops.ReadOrganizeFile{}

		if organizeText, err = readOrganizeTextOp.Run(
			u,
			os.Stdin,
			organize_text.NewMetadata(qg.RepoId),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = u.LockAndCommitOrganizeResults(
			organize_text.OrganizeResults{
				Before:     createOrganizeFileResults,
				After:      organizeText,
				Original:   skus,
				QueryGroup: qg,
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case organize_text_mode.ModeOutputOnly:
		ui.Log().Print("generate organize file and write to stdout")
		if _, err = createOrganizeFileOp.RunAndWrite(os.Stdout); err != nil {
			err = errors.Wrap(err)
			return
		}

	case organize_text_mode.ModeInteractive:
		ui.Log().Print(
			"generate temp file, write organize, open vim to edit, commit results",
		)
		var createOrganizeFileResults *organize_text.Text

		var f *os.File

		if f, err = u.GetDirectoryLayout().TempLocal.FileTempWithTemplate(
			"*." + u.GetConfig().GetFileExtensions().GetFileExtensionOrganize(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, f)

		if createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(
			f,
		); err != nil {
			err = errors.Wrapf(err, "Organize File: %q", f.Name())
			return
		}

		var organizeText *organize_text.Text

		if organizeText, err = c.readFromVim(
			u,
			f.Name(),
			createOrganizeFileResults,
			qg,
		); err != nil {
			err = errors.Wrapf(err, "Organize File: %q", f.Name())
			return
		}

		if _, err = u.LockAndCommitOrganizeResults(
			organize_text.OrganizeResults{
				Before:     createOrganizeFileResults,
				After:      organizeText,
				Original:   skus,
				QueryGroup: qg,
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = errors.Errorf("unknown mode")
		return
	}

	return
}

func (c Organize) readFromVim(
	u *env.Local,
	f string,
	results *organize_text.Text,
	qg *query.Group,
) (ot *organize_text.Text, err error) {
	openVimOp := user_ops.OpenEditor{
		VimOptions: vim_cli_options_builder.New().
			WithFileType("zit-organize").
			Build(),
	}

	if err = openVimOp.Run(u, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	readOrganizeTextOp := user_ops.ReadOrganizeFile{}

	if ot, err = readOrganizeTextOp.RunWithPath(u, f, qg.RepoId); err != nil {
		if c.handleReadChangesError(err) {
			err = nil
			ot, err = c.readFromVim(u, f, results, qg)
		} else {
			ui.Err().Printf("aborting organize")
			return
		}
	}

	return
}

func (c Organize) handleReadChangesError(err error) (tryAgain bool) {
	var errorRead organize_text.ErrorRead

	if err != nil && !errors.As(err, &errorRead) {
		ui.Err().Printf("unrecoverable organize read failure: %s", err)
		tryAgain = false
		return
	}

	tryAgain = ui.Retry("reading changes failed", err)

	return
}
