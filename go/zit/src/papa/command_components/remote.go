package command_components

import (
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/oscar/repo_remote"
)

type Remote struct {
	RemoteType repo.RemoteType
	remote     repo.WorkingCopy
}

func (cmd *Remote) SetFlagSet(f *flag.FlagSet) {
	f.Var(&cmd.RemoteType, "remote-type", fmt.Sprintf("%s", repo.GetAllRemoteTypes()))
}

func (c Remote) MakeRemote(
	env *env.Env,
	remoteArg string,
) (remote repo.WorkingCopy) {
	var err error

	switch c.RemoteType {
	case repo.RemoteTypeNativeDotenvXDG:
		if remote, err = repo_local_working_copy.MakeFromConfigAndXDGDotenvPath(
			env.Context,
			env.GetCLIConfig(),
			remoteArg,
			env.GetOptions(),
		); err != nil {
			env.CancelWithErrorAndFormat(
				err,
				"RemoteType: %q, Remote: %q",
				c.RemoteType,
				remoteArg,
			)
		}

	case repo.RemoteTypeStdioLocal:
		if remote, err = repo_remote.MakeRemoteStdioLocal(
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

	case repo.RemoteTypeStdioSSH:
		if remote, err = repo_remote.MakeRemoteStdioSSH(
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
			env.GetOptions(),
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
