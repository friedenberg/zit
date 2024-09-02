package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/organize_text"
)

func (e *Env) CommitOrganizeResults(
	results organize_text.OrganizeResults,
) (changeResults organize_text.Changes, err error) {
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

	return
}

func (e *Env) QueryGroupFromRemainingOrganizeResults(
	results organize_text.OrganizeResults,
	repoId ids.RepoId,
) (qg *query.Group, changeResults organize_text.Changes, err error) {
	if changeResults, err = organize_text.ChangesFromResults(
		e.GetConfig().PrintOptions,
		results,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	b := e.MakeQueryBuilder(
		ids.MakeGenre(genres.TrueGenre()...),
	).WithExternalLike(
		changeResults.After.AsExternalLikeSet(),
	)

	if qg, err = b.BuildQueryGroup(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *Env) CommitRemainingOrganizeResults(
	results organize_text.OrganizeResults,
	repoId ids.RepoId,
	shouldDelete bool,
) (changeResults organize_text.Changes, err error) {
	if changeResults, err = organize_text.ChangesFromResults(
		e.GetConfig().PrintOptions,
		results,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO move to bestandsaufnahme flush
	// if cs.GetAddedUnnamed().Len() == 0 && cs.GetAddedNamed().Len() == 0 {
	// 	errors.Err().Print("no changes")
	// 	return
	// }

	onChanged := func(changed sku.ExternalLike) (err error) {
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
	}

	if shouldDelete {
		withoutDelete := onChanged

		onChanged = func(changed sku.ExternalLike) (err error) {
			if err = withoutDelete(changed); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = e.GetStore().DeleteExternalLike(
				repoId,
				changed,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	if err = changeResults.After.Each(
		onChanged,
	); err != nil {
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
