package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
)

func (e *Local) MakeOrganizeOptionsWithQueryGroup(
	organizeFlags organize_text.Flags,
	qg *query.Group,
) organize_text.Options {
	return organizeFlags.GetOptions(
		e.GetConfig().PrintOptions,
		qg,
		e.SkuFormatBoxCheckedOutNoColor(),
		e.GetStore().GetAbbrStore().GetAbbr(),
		e.GetExternalLikePoolForRepoId(qg.RepoId),
	)
}

func (e *Local) LockAndCommitOrganizeResults(
	results organize_text.OrganizeResults,
) (changeResults organize_text.Changes, err error) {
	if err = e.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if changeResults, err = organize_text.ChangesFromResults(
		e.GetConfig().PrintOptions,
		results,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// ui.Debug().Print(changeResults)

	if err = changeResults.Changed.Each(
		func(changed sku.SkuType) (err error) {
			if err = e.GetStore().CreateOrUpdate(
				changed.GetSkuExternal(),
				object_mode.Make(
					object_mode.ModeMergeCheckedOut,
				),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *Local) ApplyToOrganizeOptions(oo *organize_text.Options) {
	oo.Config = e.GetConfig()
	oo.Abbr = e.GetStore().GetAbbrStore().GetAbbr()

	if !e.GetConfig().DryRun {
		return
	}

	oo.AddPrototypeAndOption(
		"dry-run",
		&organize_text.OptionCommentDryRun{MutableConfigDryRun: e.GetConfig()},
	)
}
