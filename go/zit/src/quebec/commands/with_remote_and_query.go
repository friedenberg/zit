package commands

import (
	"flag"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type CommandWithRemoteAndQuery interface {
	RunWithRemoteAndQuery(
		local *repo_local.Repo,
		remote repo.Repo,
		qg *query.Group,
	)
}

type commandWithRemoteAndQuery struct {
	command_components.Remote
	command_components.QueryGroup
	CommandWithRemoteAndQuery
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
	var qg *query.Group

	if qg, err = b.BuildQueryGroupWithRepoId(
		c.RepoId,
		c.ExternalQueryOptions,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetStore().QueryTransacted(
		qg,
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

	qg := c.MakeQueryGroup(
		c.CommandWithRemoteAndQuery,
		local,
		args[1:]...,
	)

	defer local.PrintMatchedDormantIfNecessary()

	remote := c.MakeRemote(local.Env, args[0])

	c.RunWithRemoteAndQuery(local, remote, qg)
}
