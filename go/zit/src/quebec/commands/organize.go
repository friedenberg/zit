package commands

import (
	"flag"
	"fmt"
	"os"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/organize_text_mode"
	"code.linenisgreat.com/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/src/charlie/files"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/query"
	"code.linenisgreat.com/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/src/november/umwelt"
	"code.linenisgreat.com/zit/src/papa/user_ops"
)

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
		WithDefaultSigil(kennung.SigilSchwanzen).
		WithDefaultGattungen(kennung.MakeGattung(gattung.Zettel)).
		WithRequireNonEmptyQuery()
}

func (c *Organize) CompletionGattung() kennung.Gattung {
	return kennung.MakeGattung(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
	)
}

func (c *Organize) RunWithQuery(
	u *umwelt.Umwelt,
	ms *query.Group,
) (err error) {
	u.ApplyToOrganizeOptions(&c.Options)

	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Umwelt: u,
		Options: c.GetOptions(
			u.GetKonfig().PrintOptions,
			ms,
			u.SkuFmtOrganize(),
			u.GetStore().GetAbbrStore().GetAbbr(),
		),
	}

	typen := ms.GetTypen()

	if typen.Len() == 1 {
		createOrganizeFileOp.Typ = typen.Any()
	}

	getResults := sku.MakeTransactedMutableSet()

	if err = u.GetStore().QueryWithCwd(
		ms,
		iter.MakeAddClonePoolPtrFunc(
			getResults,
			sku.GetTransactedPool(),
			sku.TransactedResetter,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	createOrganizeFileOp.Transacted = getResults

	switch c.Mode {
	case organize_text_mode.ModeCommitDirectly:
		ui.Log().Print("neither stdin or stdout is a tty")
		ui.Log().Print("generate organize, read from stdin, commit")

		var createOrganizeFileResults *organize_text.Text

		var f *os.File

		if f, err = files.TempFileWithPattern(
			"*." + u.GetKonfig().FileExtensions.Organize,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(
			f,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		var ot2 *organize_text.Text

		readOrganizeTextOp := user_ops.ReadOrganizeFile{}

		if ot2, err = readOrganizeTextOp.Run(u, os.Stdin); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = u.Lock(); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, u.Unlock)

		commitOrganizeTextOp := user_ops.CommitOrganizeFile{
			Umwelt: u,
		}

		if _, err = commitOrganizeTextOp.Run(
			u,
			createOrganizeFileResults,
			ot2,
			getResults,
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

		if f, err = u.Standort().FileTempLocalWithTemplate(
			"*." + u.GetKonfig().FileExtensions.Organize,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if createOrganizeFileResults, err = createOrganizeFileOp.RunAndWrite(
			f,
		); err != nil {
			err = errors.Wrapf(err, "Organize File: %q", f.Name())
			return
		}

		var ot2 *organize_text.Text

		if ot2, err = c.readFromVim(u, f.Name(), createOrganizeFileResults, ms); err != nil {
			err = errors.Wrapf(err, "Organize File: %q", f.Name())
			return
		}

		if err = u.Lock(); err != nil {
			err = errors.Wrapf(err, "Organize File: %q", f.Name())
			return
		}

		defer errors.Deferred(&err, u.Unlock)

		commitOrganizeTextOp := user_ops.CommitOrganizeFile{
			Umwelt: u,
		}

		if _, err = commitOrganizeTextOp.Run(
			u,
			createOrganizeFileResults,
			ot2,
			getResults,
		); err != nil {
			err = errors.Wrapf(err, "Organize File: %q", f.Name())
			return
		}

	default:
		err = errors.Errorf("unknown mode")
		return
	}

	return
}

func (c Organize) readFromVim(
	u *umwelt.Umwelt,
	f string,
	results *organize_text.Text,
	q *query.Group,
) (ot *organize_text.Text, err error) {
	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithFileType("zit-organize").
			Build(),
	}

	if _, err = openVimOp.Run(u, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	readOrganizeTextOp := user_ops.ReadOrganizeFile{}

	if ot, err = readOrganizeTextOp.RunWithPath(u, f); err != nil {
		if c.handleReadChangesError(err) {
			err = nil
			ot, err = c.readFromVim(u, f, results, q)
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

	ui.Err().Printf("reading changes failed: %q", err)
	ui.Err().Printf("would you like to edit and try again? (y/*)")

	var answer rune
	var n int

	if n, err = fmt.Scanf("%c", &answer); err != nil {
		tryAgain = false
		ui.Err().Printf("failed to read answer: %s", err)
		return
	}

	if n != 1 {
		tryAgain = false
		ui.Err().Printf("failed to read at exactly 1 answer: %s", err)
		return
	}

	if answer == 'y' || answer == 'Y' {
		tryAgain = true
	}

	return
}
