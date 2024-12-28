package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
	"code.linenisgreat.com/zit/go/zit/src/oscar/repo_remote"
)

type Remote struct {
	RemoteType repo.RemoteType
	remote     repo.Repo
}

func (cmd *Remote) SetFlagSet(f *flag.FlagSet) {
	f.Var(&cmd.RemoteType, "remote-type", "TODO")
}

func (c Remote) MakeRemote(
	env *env.Env,
	remoteArg string,
) (remote repo.Repo) {
	var err error

	switch c.RemoteType {
	case repo.RemoteTypeNativeDotenvXDG:
		if remote, err = repo_local.MakeFromConfigAndXDGDotenvPath(
			env.Context,
			env.GetCLIConfig(),
			remoteArg,
		); err != nil {
			env.CancelWithErrorAndFormat(
				err,
				"RemoteType: %q, Remote: %q",
				c.RemoteType,
				remoteArg,
			)
		}

	case repo.RemoteTypeStdioLocal:
		if remote, err = repo_remote.MakeRemoteStdio(
			env,
			remoteArg,
		); err != nil {
			env.CancelWithErrorAndFormat(
				err,
				"RemoteType: %q, Remote: %q",
				c.RemoteType,
				remoteArg,
			)
		}

	case repo.RemoteTypeSocketUnix:
		if remote, err = repo_remote.MakeRemoteHTTPFromXDGDotenvPath(
			env.Context,
			env.GetCLIConfig(),
			remoteArg,
		); err != nil {
			env.CancelWithErrorAndFormat(
				err,
				"RemoteType: %q, Remote: %q",
				c.RemoteType,
				remoteArg,
			)
		}

	default:
		env.CancelWithNotImplemented()
	}

	return
}
