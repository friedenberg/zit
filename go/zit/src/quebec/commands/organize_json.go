package commands

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type OrganizeJSON struct {
	Or bool
	organize_text.Flags

	Filter script_value.ScriptValue
}

func init() {
	registerCommandWithQuery(
		"organize-json",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &OrganizeJSON{
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

			c.AddToFlagSet(f)

			return c
		},
	)
}

func (c *OrganizeJSON) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung(
		gattung.Zettel,
	)
}

func (c *OrganizeJSON) CompletionGattung() kennung.Gattung {
	return kennung.MakeGattung(
		gattung.Zettel,
		gattung.Etikett,
		gattung.Typ,
	)
}

func (c *OrganizeJSON) RunWithQuery(
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

	u.GetKonfig().DryRun = true
	u.GetKonfig().PrintOptions.Abbreviations.Hinweisen = false

	var transacted []sku_fmt.Json

	dec := json.NewDecoder(bufio.NewReader(u.In()))

	if err = dec.Decode(&transacted); err != nil {
		err = errors.Wrap(err)
		return
	}

	getResults := sku.MakeTransactedMutableSetKennung()

	for _, j := range transacted {
		sk := sku.GetTransactedPool().Get()

		if err = j.ToTransacted(sk, u.Standort()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = getResults.Add(sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	createOrganizeFileOp.Transacted = getResults

	ui.Log().Print(
		"generate temp file, write organize, open vim to edit, commit results",
	)
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
		Umwelt:     u,
		OutputJSON: true,
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

	return
}

func (c OrganizeJSON) readFromVim(
	u *umwelt.Umwelt,
	p string,
	results *organize_text.Text,
	q *query.Group,
) (ot *organize_text.Text, err error) {
	openVimOp := user_ops.OpenVim{
		Options: vim_cli_options_builder.New().
			WithFileType("zit-organize").
			Build(),
	}

	if err = openVimOp.Run(u, p); err != nil {
		err = errors.Wrap(err)
		return
	}

	readOrganizeTextOp := user_ops.ReadOrganizeFile{}

	if ot, err = readOrganizeTextOp.RunWithPath(u, p); err != nil {
		if c.handleReadChangesError(err) {
			err = nil
			ot, err = c.readFromVim(u, p, results, q)
		} else {
			ui.Err().Printf("aborting organize")
			return
		}
	}

	return
}

func (c OrganizeJSON) handleReadChangesError(err error) (tryAgain bool) {
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
