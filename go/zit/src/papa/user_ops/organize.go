package user_ops

import (
	"fmt"
	"os"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/src/charlie/files"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/objekte_collections"
	"code.linenisgreat.com/zit/src/juliett/query"
	"code.linenisgreat.com/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/src/november/umwelt"
)

type Organize struct {
	*umwelt.Umwelt
	metadatei.Metadatei
}

func (u Organize) Run(qg *query.Group, skus sku.TransactedSet) (err error) {
	if qg == nil {
		b := u.MakeQueryBuilder(
			kennung.MakeGattung(gattung.TrueGattung()...),
		).WithTransacted(
			skus,
		)

		if qg, err = b.BuildQueryGroup(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	otFlags := organize_text.MakeFlagsWithMetadatei(u.Metadatei)
	u.ApplyToOrganizeOptions(&otFlags.Options)
	// otFlags.Abbr = u.StoreObjekten().GetAbbrStore().AbbreviateHinweis
	mwk := objekte_collections.MakeMutableSetMetadateiWithKennung()
	skus.Each(
		func(z *sku.Transacted) (err error) {
			return mwk.Add(z)
		},
	)
	otFlags.Transacted = mwk

	createOrganizeFileOp := CreateOrganizeFile{
		Umwelt: u.Umwelt,
		Options: otFlags.GetOptions(
			u.Konfig().PrintOptions,
			qg,
			u.SkuFormatOldOrganize(),
			u.SkuFmtNewOrganize(),
			u.GetStore().GetAbbrStore().GetAbbr(),
		),
	}

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

	for {
		openVimOp := OpenVim{
			Options: vim_cli_options_builder.New().
				WithFileType("zit-organize").
				Build(),
		}

		if _, err = openVimOp.Run(u.Umwelt, f.Name()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = u.Reset(); err != nil {
			err = errors.Wrap(err)
			return
		}

		readOrganizeTextOp := ReadOrganizeFile{
			Umwelt: u.Umwelt,
		}

		if ot2, err = readOrganizeTextOp.RunWithFile(f.Name(), qg); err != nil {
			if u.handleReadChangesError(err) {
				err = nil
				continue
			} else {
				errors.Err().Printf("aborting organize")
				return
			}
		}

		break
	}

	commitOrganizeTextOp := CommitOrganizeFile{
		Umwelt: u.Umwelt,
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	if _, err = commitOrganizeTextOp.Run(
		u.Umwelt,
		createOrganizeFileResults,
		ot2,
		mwk,
	); err != nil {
		err = errors.Wrap(err)
		return
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
