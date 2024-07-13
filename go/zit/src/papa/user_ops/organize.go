package user_ops

import (
	"fmt"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type Organize struct {
	*umwelt.Umwelt
	object_metadata.Metadata
}

func (u Organize) Run(qg *query.Group, skus sku.TransactedSet) (err error) {
	if qg == nil {
		b := u.MakeQueryBuilder(
			ids.MakeGenre(genres.TrueGenre()...),
		).WithTransacted(
			skus,
		)

		if qg, err = b.BuildQueryGroup(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	otFlags := organize_text.MakeFlagsWithMetadatei(u.Metadata)
	u.ApplyToOrganizeOptions(&otFlags.Options)
	// otFlags.Abbr = u.StoreObjekten().GetAbbrStore().AbbreviateHinweis
	mwk := sku.MakeTransactedMutableSet()
	skus.Each(
		func(z *sku.Transacted) (err error) {
			return mwk.Add(z)
		},
	)
	otFlags.Transacted = mwk

	createOrganizeFileOp := CreateOrganizeFile{
		Umwelt: u.Umwelt,
		Options: otFlags.GetOptions(
			u.GetKonfig().PrintOptions,
			qg,
			u.SkuFmtOrganize(),
			u.GetStore().GetAbbrStore().GetAbbr(),
		),
	}

	var createOrganizeFileResults *organize_text.Text

	var f *os.File

	if f, err = u.Standort().FileTempLocalWithTemplate(
		"*." + u.GetKonfig().FileExtensions.Organize,
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

	var ot2 *organize_text.Text

	for {
		openVimOp := OpenVim{
			Options: vim_cli_options_builder.New().
				WithFileType("zit-organize").
				Build(),
		}

		if err = openVimOp.Run(u.Umwelt, f.Name()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = u.Reset(); err != nil {
			err = errors.Wrap(err)
			return
		}

		readOrganizeTextOp := ReadOrganizeFile{}

		if _, err = f.Seek(0, io.SeekStart); err != nil {
			err = errors.Wrap(err)
			return
		}

		if ot2, err = readOrganizeTextOp.Run(u.Umwelt, f); err != nil {
			if u.handleReadChangesError(err) {
				err = nil
				continue
			} else {
				ui.Err().Printf("aborting organize")
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
