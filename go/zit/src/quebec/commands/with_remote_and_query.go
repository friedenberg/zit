package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type WithRemoteAndQuery interface {
	Run(
		local *local_working_copy.Repo,
		remote repo.Archive,
		qg *query.Group,
		options repo.RemoteTransferOptions,
	)
}

type commandWithRemoteAndQuery struct {
	command_components.RemoteTransfer
	command_components.QueryGroup
	Command WithRemoteAndQuery
}

func (cmd *commandWithRemoteAndQuery) SetFlagSet(f *flag.FlagSet) {
	cmd.RemoteTransfer.SetFlagSet(f)
	cmd.QueryGroup.SetFlagSet(f)

	if cwf, ok := cmd.Command.(interfaces.CommandComponent); ok {
		cwf.SetFlagSet(f)
	}
}

func (c commandWithRemoteAndQuery) CompleteWithRepo(
	u *local_working_copy.Repo,
	args ...string,
) (err error) {
	c.QueryGroup.CompleteWithRepo(
		c.Command,
		u,
		args...,
	)

	return
}

func (c commandWithRemoteAndQuery) Run(
	local *local_working_copy.Repo,
	args ...string,
) {
	if len(args) < 1 {
		// TODO add info about remote options
		local.CancelWithBadRequestf("requires a remote to be specified")
	}

	qg := c.MakeQueryGroup(
		query.MakeBuilderOptions(c.Command),
		local,
		args[1:]...,
	)

	remote := c.MakeWorkingCopy(local.Env, args[0])

	c.Command.Run(local, remote, qg, c.RemoteTransferOptions)
}
