package commands

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type CommandWithRemoteAndQuery interface {
	RunWithRemoteAndQuery(
		local *env.Local,
		remote env.Env,
		qg *query.Group,
	)
}

type commandWithRemoteAndQuery struct {
	TheirXDGDotenv string
	UseSocket      bool

	CommandWithRemoteAndQuery

	remote env.Env
	sku.ExternalQueryOptions
	*query.Group
}

func (c commandWithRemoteAndQuery) Complete(
	u *env.Local,
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

func (c commandWithRemoteAndQuery) Run(
	local *env.Local,
	args ...string,
) {
	if len(args) < 1 && c.TheirXDGDotenv == "" {
		// TODO add info about remote options
		local.CancelWithError(errors.BadRequestf("Pulling requires a remote to be specified"))
		return
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
			return
		}
	}

	var remote env.Env

	{
		var err error

		if c.TheirXDGDotenv != "" {
			if c.UseSocket {
				if remote, err = env.MakeRemoteHTTPFromXDGDotenvPath(
					local.Context,
					local.GetConfig(),
					c.TheirXDGDotenv,
				); err != nil {
					local.CancelWithError(err)
					return
				}
			} else {
				if remote, err = env.MakeLocalFromConfigAndXDGDotenvPath(
					local.Context,
					local.GetConfig(),
					c.TheirXDGDotenv,
				); err != nil {
					local.CancelWithError(err)
					return
				}
			}
		} else {
			local.CancelWithError(todo.Implement())
			return
		}
	}

	c.RunWithRemoteAndQuery(local, remote, c.Group)

	return
}
