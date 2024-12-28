package commands

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type commandWithQuery struct {
	CommandWithQuery
	command_components.QueryGroup
}

type CompletionGenresGetter interface {
	CompletionGenres() ids.Genre
}

func (c commandWithQuery) CompleteWithRepo(
	u *repo_local.Repo,
	args ...string,
) {
	var cgg CompletionGenresGetter
	ok := false

	if cgg, ok = c.CommandWithQuery.(CompletionGenresGetter); !ok {
		return
	}

	w := sku_fmt.MakeWriterComplete(os.Stdout)
	defer u.MustClose(w)

	b := u.MakeQueryBuilderExcludingHidden(cgg.CompletionGenres())

	var qg *query.Group

	{
		var err error

		if qg, err = b.BuildQueryGroupWithRepoId(
			c.RepoId,
			c.ExternalQueryOptions,
		); err != nil {
			u.Context.CancelWithError(err)
		}
	}

	if err := u.GetStore().QueryTransacted(
		qg,
		w.WriteOneTransacted,
	); err != nil {
		u.Context.CancelWithError(err)
	}
}

func (c commandWithQuery) RunWithRepo(
	u *repo_local.Repo,
	args ...string,
) {
	var qg *query.Group

	{
		var err error

		if qg, err = u.MakeQueryGroup(
			c.CommandWithQuery,
			c.RepoId,
			c.ExternalQueryOptions,
			args...,
		); err != nil {
			u.CancelWithError(err)
		}
	}

	defer u.PrintMatchedDormantIfNecessary()

	c.RunWithQuery(u, qg)
}
