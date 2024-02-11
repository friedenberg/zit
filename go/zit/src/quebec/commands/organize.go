package commands

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit-go/src/bravo/files"
	"code.linenisgreat.com/zit-go/src/charlie/gattung"
	"code.linenisgreat.com/zit-go/src/charlie/script_value"
	"code.linenisgreat.com/zit-go/src/delta/gattungen"
	"code.linenisgreat.com/zit-go/src/echo/kennung"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
	"code.linenisgreat.com/zit-go/src/india/matcher"
	"code.linenisgreat.com/zit-go/src/india/objekte_collections"
	"code.linenisgreat.com/zit-go/src/kilo/organize_text"
	"code.linenisgreat.com/zit-go/src/oscar/umwelt"
	"code.linenisgreat.com/zit-go/src/papa/user_ops"
)

type Organize struct {
	Or bool
	organize_text.Flags
	Mode organize_text.Mode

	Filter script_value.ScriptValue
}

func init() {
	registerCommandWithQuery(
		"organize",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Organize{
				Flags: organize_text.MakeFlags(),
			}

			f.BoolVar(
				&c.Or,
				"or",
				false,
				"allow optional criteria instead of required",
			)
			f.Var(
				&c.Filter,
				"filter",
				"a script to run for each file to transform it the standard zettel format",
			)
			f.Var(&c.Mode, "mode", "mode used for handling stdin and stdout")

			c.Flags.AddToFlagSet(f)

			return c
		},
	)
}

func (c *Organize) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Zettel,
	)
}

func (c *Organize) CompletionGattung() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
	)
}

func (c *Organize) RunWithQuery(
	u *umwelt.Umwelt,
	ms matcher.Query,
) (err error) {
	u.ApplyToOrganizeOptions(&c.Options)

	createOrganizeFileOp := user_ops.CreateOrganizeFile{
		Umwelt: u,
		Options: c.GetOptions(
			u.Konfig().PrintOptions,
			ms,
			u.SkuFormatOldOrganize(),
			u.SkuFmtNewOrganize(),
		),
	}

	typen := ms.GetTypen()

	switch typen.Len() {
	case 0:
		break

	case 1:
		createOrganizeFileOp.Typ = typen.Any()

	default:
		err = errors.Errorf(
			"only one typ is supported for organize, but got %q",
			typen,
		)
		return
	}

	var l sync.Mutex
	getResults := objekte_collections.MakeMutableSetMetadateiWithKennung()

	if err = u.StoreObjekten().QueryWithCwd(
		ms,
		func(tl *sku.Transacted) (err error) {
			mwk := sku.GetTransactedPool().Get()

			if err = mwk.SetFromSkuLike(tl); err != nil {
				err = errors.Wrap(err)
				return
			}

			// TODO-P1 determine if this is necessary
			var h *kennung.Hinweis
			h = &kennung.Hinweis{}

			if err = h.Set(mwk.GetKennung().String()); err == nil {
				if h, err = u.StoreObjekten().GetAbbrStore().Hinweis().ExpandString(
					h.String(),
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = mwk.SetKennungLike(h); err != nil {
					err = errors.Wrap(err)
					return
				}
			} else {
				err = nil
			}

			l.Lock()
			defer l.Unlock()

			return getResults.Add(mwk)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	createOrganizeFileOp.Transacted = getResults

	switch c.Mode {
	case organize_text.ModeCommitDirectly:
		errors.Log().Print("neither stdin or stdout is a tty")
		errors.Log().Print("generate organize, read from stdin, commit")

		var createOrganizeFileResults *organize_text.Text

		var f *os.File

		if f, err = files.TempFileWithPattern(
			"*." + u.Konfig().FileExtensions.Organize,
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

		readOrganizeTextOp := user_ops.ReadOrganizeFile{
			Umwelt: u,
			Reader: os.Stdin,
		}

		if ot2, err = readOrganizeTextOp.Run(ms); err != nil {
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
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case organize_text.ModeOutputOnly:
		errors.Log().Print("generate organize file and write to stdout")
		if _, err = createOrganizeFileOp.RunAndWrite(os.Stdout); err != nil {
			err = errors.Wrap(err)
			return
		}

	case organize_text.ModeInteractive:
		errors.Log().Print(
			"generate temp file, write organize, open vim to edit, commit results",
		)
		var createOrganizeFileResults *organize_text.Text

		var f *os.File

		if f, err = files.TempFileWithPattern(
			"*." + u.Konfig().FileExtensions.Organize,
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

		if ot2, err = c.readFromVim(u, f.Name(), createOrganizeFileResults, ms); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = u.Reset(); err != nil {
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
	u *umwelt.Umwelt,
	f string,
	results *organize_text.Text,
	q matcher.Query,
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

	readOrganizeTextOp := user_ops.ReadOrganizeFile{
		Umwelt: u,
	}

	if ot, err = readOrganizeTextOp.RunWithFile(f, q); err != nil {
		if c.handleReadChangesError(err) {
			err = nil
			ot, err = c.readFromVim(u, f, results, q)
		} else {
			errors.Err().Printf("aborting organize")
			return
		}
	}

	return
}

func (c Organize) handleReadChangesError(err error) (tryAgain bool) {
	var errorRead organize_text.ErrorRead

	if err != nil && !errors.As(err, &errorRead) {
		errors.Err().Printf("unrecoverable organize read failure: %s", err)
		tryAgain = false
		return
	}

	errors.Err().Printf("reading changes failed: %q", err)
	errors.Err().Printf("would you like to edit and try again? (y/*)")

	var answer rune
	var n int

	if n, err = fmt.Scanf("%c", &answer); err != nil {
		tryAgain = false
		errors.Err().Printf("failed to read answer: %s", err)
		return
	}

	if n != 1 {
		tryAgain = false
		errors.Err().Printf("failed to read at exactly 1 answer: %s", err)
		return
	}

	if answer == 'y' || answer == 'Y' {
		tryAgain = true
	}

	return
}
