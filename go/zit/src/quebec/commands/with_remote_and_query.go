package commands

import (
	"flag"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
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
	RemoteType repo.RemoteType

	CommandWithRemoteAndQuery

	remote repo.Repo
	sku.ExternalQueryOptions
	*query.Group
}

func (cmd *commandWithRemoteAndQuery) SetFlagSet(f *flag.FlagSet) {
	f.Var(&cmd.RepoId, "kasten", "none or Browser")
	f.Var(&cmd.RemoteType, "remote-type", "TODO")
	f.BoolVar(&cmd.ExcludeUntracked, "exclude-untracked", false, "")
	f.BoolVar(&cmd.ExcludeRecognized, "exclude-recognized", false, "")

	if cwf, ok := cmd.CommandWithRemoteAndQuery.(CommandWithFlags); ok {
		cwf.SetFlagSet(f)
	}
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
	if len(args) < 1 {
		// TODO add info about remote options
		local.CancelWithBadRequestf("requires a remote to be specified")
	}

	{
		var err error

		if c.Group, err = local.MakeQueryGroup(
			c.CommandWithRemoteAndQuery,
			c.RepoId,
			c.ExternalQueryOptions,
			args[1:]...,
		); err != nil {
			local.CancelWithError(err)
		}
	}

	defer local.PrintMatchedDormantIfNecessary()

	remote := c.makeRemote(local, args[0])

	c.RunWithRemoteAndQuery(local, remote, c.Group)
}

func (c commandWithRemoteAndQuery) makeRemote(
	local *repo_local.Repo,
	remoteArg string,
) (remote repo.Repo) {
	var err error

	switch c.RemoteType {
	case repo.RemoteTypeNativeDotenvXDG:
		if remote, err = repo_local.MakeFromConfigAndXDGDotenvPath(
			local.Context,
			local.GetConfig(),
			remoteArg,
		); err != nil {
			local.CancelWithErrorAndFormat(
				err,
				"RemoteType: %q, Remote: %q",
				c.RemoteType,
				remoteArg,
			)
		}

	case repo.RemoteTypeSocketUnix:
		if remote, err = repo_remote.MakeRemoteHTTPFromXDGDotenvPath(
			local.Context,
			local.GetConfig(),
			remoteArg,
		); err != nil {
			local.CancelWithErrorAndFormat(
				err,
				"RemoteType: %q, Remote: %q",
				c.RemoteType,
				remoteArg,
			)
		}

	default:
		local.CancelWithNotImplemented()
	}

	return
}
