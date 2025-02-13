package local_working_copy

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/organize_text"
)

func (e *Repo) MakeOrganizeOptionsWithQueryGroup(
	organizeFlags organize_text.Flags,
	qg *query.Query,
) organize_text.Options {
	return organizeFlags.GetOptions(
		e.GetConfig().GetCLIConfig().PrintOptions,
		query.GetTags(qg),
		e.SkuFormatBoxCheckedOutNoColor(),
		e.GetStore().GetAbbrStore().GetAbbr(),
		e.GetExternalLikePoolForRepoId(qg.RepoId),
	)
}

func (repo *Repo) LockAndCommitOrganizeResults(
	results organize_text.OrganizeResults,
) (changeResults organize_text.Changes, err error) {
	if changeResults, err = organize_text.ChangesFromResults(
		repo.GetConfig().GetCLIConfig().PrintOptions,
		results,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	repo.Must(repo.Lock)

	count := changeResults.Changed.Len()

	if count > 30 {
		if !repo.Confirm(
			fmt.Sprintf(
				"a large number (%d) of objects are being changed. continue to commit?",
				count,
			),
		) {
			// TODO output organize file used
			repo.CancelWithBadRequestf("aborting")
			return
		}
	}

	for _, changed := range changeResults.Changed.AllSkuAndIndex() {
		if err = repo.GetStore().CreateOrUpdate(
			changed.GetSkuExternal(),
			sku.StoreOptions{
				MergeCheckedOut: true,
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	repo.Must(repo.Unlock)

	return
}

func (e *Repo) ApplyToOrganizeOptions(oo *organize_text.Options) {
	oo.Config = e.GetConfig()
	oo.Abbr = e.GetStore().GetAbbrStore().GetAbbr()

	if !e.GetConfig().GetCLIConfig().IsDryRun() {
		return
	}

	oo.AddPrototypeAndOption(
		"dry-run",
		&organize_text.OptionCommentDryRun{
			MutableConfigDryRun: e.GetConfig(),
		},
	)
}
