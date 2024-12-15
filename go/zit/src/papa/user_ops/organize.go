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
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Organize struct {
	*env.Local
	organize_text.Metadata
	DontUseQueryGroupForOrganizeMetadata bool
}

func (op Organize) RunWithQueryGroup(
	qg *query.Group,
) (organizeResults organize_text.OrganizeResults, err error) {
	skus := sku.MakeSkuTypeSetMutable()
	var l sync.RWMutex

	if err = op.GetStore().QueryTransactedAsSkuType(
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

	if organizeResults, err = op.RunWithExternalLike(qg, skus); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO remove
func (op Organize) RunWithTransacted(
	qg *query.Group,
	transacted sku.TransactedSet,
) (organizeResults organize_text.OrganizeResults, err error) {
	skus := sku.MakeSkuTypeSetMutable()

	for z := range transacted.All() {
		clone := sku.CloneSkuTypeFromTransacted(
			z.GetSku(),
			checked_out_state.Unknown,
		)

		skus.Add(clone)
	}

	if organizeResults, err = op.RunWithExternalLike(qg, skus); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op Organize) RunWithExternalLike(
	qg *query.Group,
	skus sku.SkuTypeSet,
) (organizeResults organize_text.OrganizeResults, err error) {
	organizeResults.Original = skus
	organizeResults.QueryGroup = qg

	var repoId ids.RepoId

	if qg != nil {
		repoId = qg.RepoId
	}

	if organizeResults.QueryGroup == nil ||
		op.DontUseQueryGroupForOrganizeMetadata {
		b := op.MakeQueryBuilder(
			ids.MakeGenre(genres.TrueGenre()...),
		).WithExternalLike(
			skus,
		)

		if organizeResults.QueryGroup, err = b.BuildQueryGroup(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	organizeResults.QueryGroup.RepoId = repoId

	organizeFlags := organize_text.MakeFlagsWithMetadata(op.Metadata)
	op.ApplyToOrganizeOptions(&organizeFlags.Options)
	organizeFlags.Skus = skus

	createOrganizeFileOp := CreateOrganizeFile{
		Local: op.Local,
		Options: op.Local.MakeOrganizeOptionsWithQueryGroup(
			organizeFlags,
			organizeResults.QueryGroup,
		),
	}

	typen := organizeResults.QueryGroup.GetTypes()

	if typen.Len() == 1 {
		createOrganizeFileOp.Type = typen.Any()
	}

	var f *os.File

	if f, err = op.GetDirectoryLayout().TempLocal.FileTempWithTemplate(
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

		if err = openVimOp.Run(op.Local, f.Name()); err != nil {
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
			op.Local,
			f,
			organize_text.NewMetadataWithOptionCommentLookup(
				organizeResults.Before.Metadata.RepoId,
				op.GetPrototypeOptionComments(),
			),
		); err != nil {
			if op.handleReadChangesError(err) {
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
