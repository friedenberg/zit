package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/organize_text"
)

func (e *Env) MakeOrganizeOptionsWithQueryGroup(
  organizeFlags organize_text.Flags,
  qg *query.Group,
) organize_text.Options {
	return organizeFlags.GetOptions(
		e.GetConfig().PrintOptions,
		qg,
		e.SkuFormatBox(),
		e.GetStore().GetAbbrStore().GetAbbr(),
		e.GetExternalLikePoolForRepoId(organizeResults.QueryGroup.RepoId),
	)
}

func (e *Env) LockAndCommitOrganizeResults(
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

	if err = changeResults.Changed.Each(
		func(changed sku.ExternalLike) (err error) {
			if err = e.GetStore().CreateOrUpdate(
				changed,
				objekte_mode.Make(
					objekte_mode.ModeMergeCheckedOut,
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

func (e *Env) ApplyToOrganizeOptions(oo *organize_text.Options) {
	oo.Config = e.GetConfig()
	oo.Abbr = e.GetStore().GetAbbrStore().GetAbbr()

	if !e.GetConfig().DryRun {
		return
	}

	oo.AddPrototypeAndOption(
		"dry-run",
		&organize_text.OptionCommentDryRun{ConfigDryRun: e.GetConfig()},
	)
}
