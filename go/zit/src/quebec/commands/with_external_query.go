package commands

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type commandWithQuery struct {
	CommandWithQuery
	sku.ExternalQueryOptions
	*query.Group
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
	defer u.Context.Closer(w)

	b := u.MakeQueryBuilderExcludingHidden(cgg.CompletionGenres())

	{
		var err error

		if c.Group, err = b.BuildQueryGroupWithRepoId(
			c.RepoId,
			c.ExternalQueryOptions,
		); err != nil {
			u.Context.CancelWithError(err)
			return
		}
	}

	if err := u.GetStore().QueryTransacted(
		c.Group,
		w.WriteOneTransacted,
	); err != nil {
		u.Context.CancelWithError(err)
		return
	}
}

func (c commandWithQuery) RunWithRepo(
	u *repo_local.Repo,
	args ...string,
) {
	{
		var err error
		if c.Group, err = u.MakeQueryGroup(
			c.CommandWithQuery,
			c.RepoId,
			c.ExternalQueryOptions,
			args...,
		); err != nil {
			u.Context.CancelWithError(err)
			return
		}
	}

	defer u.PrintMatchedDormantIfNecessary()

	if err := c.RunWithQuery(u, c.Group); err != nil {
		u.Context.CancelWithError(err)
		return
	}
}
