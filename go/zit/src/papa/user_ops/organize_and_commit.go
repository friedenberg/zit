package user_ops

import (
	"fmt"
	"io"
	"os"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/vim_cli_options_builder"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type OrganizeAndCommit struct {
	*env.Env
	organize_text.Metadata
}

func (u OrganizeAndCommit) RunWithQueryGroup(
	qg *query.Group,
) (results organize_text.OrganizeResults, err error) {
	skus := sku.MakeExternalLikeMutableSet()
	var l sync.Mutex

	if err = u.GetStore().Query(
		qg,
		func(el sku.ExternalLike) (err error) {
			l.Lock()
			defer l.Unlock()
			return skus.Add(el.Clone())
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if results, err = u.RunWithExternalLike(qg, skus); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (u OrganizeAndCommit) RunWithTransacted(
	qg *query.Group,
	transacted sku.TransactedSet,
) (results organize_text.OrganizeResults, err error) {
	skus := sku.MakeExternalLikeMutableSet()
	transacted.Each(
		func(z *sku.Transacted) (err error) {
			return skus.Add(z)
		},
	)

	if results, err = u.RunWithExternalLike(qg, skus); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (u OrganizeAndCommit) RunWithExternalLike(
	qg *query.Group,
	skus sku.ExternalLikeSet,
) (results organize_text.OrganizeResults, err error) {
	if qg == nil {
		b := u.MakeQueryBuilder(
			ids.MakeGenre(genres.TrueGenre()...),
		).WithExternalLike(
			skus,
		)

		if qg, err = b.BuildQueryGroup(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	// otFlags.Abbr = u.StoreObjekten().GetAbbrStore().AbbreviateHinweis
	organizeFlags := organize_text.MakeFlagsWithMetadata(u.Metadata)
	u.ApplyToOrganizeOptions(&organizeFlags.Options)
	organizeFlags.Skus = skus

	createOrganizeFileOp := CreateOrganizeFile{
		Env: u.Env,
		Options: organizeFlags.GetOptions(
			u.GetConfig().PrintOptions,
			qg,
			u.SkuFmtOrganize(qg.RepoId),
			u.GetStore().GetAbbrStore().GetAbbr(),
			u.GetExternalLikePoolForRepoId(qg.RepoId),
		),
	}

	typen := qg.GetTypes()

	if typen.Len() == 1 {
		createOrganizeFileOp.Type = typen.Any()
	}

	var f *os.File

	if f, err = u.GetFSHome().FileTempLocalWithTemplate(
		"*." + u.GetConfig().FileExtensions.Organize,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if results.Before, err = createOrganizeFileOp.RunAndWrite(
		f,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	for {
		openVimOp := OpenVim{
			Options: vim_cli_options_builder.New().
				WithFileType("zit-organize").
				Build(),
		}

		if err = openVimOp.Run(u.Env, f.Name()); err != nil {
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

		if results.After, err = readOrganizeTextOp.Run(
			u.Env,
			f,
			qg.RepoId,
			organize_text.NewMetadata(),
		); err != nil {
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

	results.Original = skus
	results.QueryGroup = qg

	if _, err = u.CommitOrganizeResults(results); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c OrganizeAndCommit) handleReadChangesError(err error) (tryAgain bool) {
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
