package commands

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
	"code.linenisgreat.com/zit/go/zit/src/oscar/repo_remote"
)

type CommandWithRemoteAndQuery interface {
	RunWithRemoteAndQuery(
		local *repo_local.Repo,
		remote repo.Repo,
		qg *query.Group,
	)
}

type commandWithRemoteAndQuery struct {
	TheirXDGDotenv string
	UseSocket      bool

	CommandWithRemoteAndQuery

	remote repo.Repo
	sku.ExternalQueryOptions
	*query.Group
}

func (c commandWithRemoteAndQuery) CompleteWithRepo(
	u *repo_local.Repo,
	args ...string,
) (err error) {
	var cgg CompletionGenresGetter
	ok := false

	if cgg, ok = c.CommandWithRemoteAndQuery.(CompletionGenresGetter); !ok {
		return
	}

	w := sku_fmt.MakeWriterComplete(os.Stdout)
	defer errors.DeferredCloser(&err, w)

	b := u.MakeQueryBuilderExcludingHidden(cgg.CompletionGenres())

	if c.Group, err = b.BuildQueryGroupWithRepoId(
		c.RepoId,
		c.ExternalQueryOptions,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetStore().QueryTransacted(
		c.Group,
		w.WriteOneTransacted,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c commandWithRemoteAndQuery) RunWithRepo(
	local *repo_local.Repo,
	args ...string,
) {
	if len(args) < 1 && c.TheirXDGDotenv == "" {
		// TODO add info about remote options
		local.CancelWithBadRequestf("Pulling requires a remote to be specified")
	}

	{
		var err error

		if c.Group, err = local.MakeQueryGroup(
			c.CommandWithRemoteAndQuery,
			c.RepoId,
			c.ExternalQueryOptions,
			args...,
		); err != nil {
			local.CancelWithError(err)
		}
	}

	defer local.PrintMatchedDormantIfNecessary()

	var remote repo.Repo

	{
		var err error

		if c.TheirXDGDotenv != "" {
			if c.UseSocket {
				if remote, err = repo_remote.MakeRemoteHTTPFromXDGDotenvPath(
					local.Context,
					local.GetConfig(),
					c.TheirXDGDotenv,
				); err != nil {
					local.CancelWithError(err)
				}
			} else {
				if remote, err = repo_local.MakeFromConfigAndXDGDotenvPath(
					local.Context,
					local.GetConfig(),
					c.TheirXDGDotenv,
				); err != nil {
					local.CancelWithError(err)
				}
			}
		} else {
			local.CancelWithError(todo.Implement())
		}
	}

	c.RunWithRemoteAndQuery(local, remote, c.Group)
}
