package command_components

import (
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/oscar/repo_remote"
)

type Remote struct {
	RemoteType repo.RemoteType
}

func (cmd *Remote) SetFlagSet(f *flag.FlagSet) {
	f.Var(&cmd.RemoteType, "remote-type", fmt.Sprintf("%s", repo.GetAllRemoteTypes()))
}

func (cmd Remote) MakeWorkingCopyFromFlagSet(
	env *env.Env,
	f *flag.FlagSet,
) (remote repo.WorkingCopy) {
	if len(f.Args()) == 0 {
		// TODO add info about remote options
		env.CancelWithBadRequestf("requires a remote to be specified")
	}

	return cmd.MakeWorkingCopy(env, f.Args()[0])
}

func (c Remote) MakeArchive(
	env *env.Env,
	remoteArg string,
) (remote repo.Archive) {
	var err error

	switch c.RemoteType {
	case repo.RemoteTypeNativeDotenvXDG:
		if remote, err = local_working_copy.MakeFromConfigAndXDGDotenvPath(
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

func (c Remote) MakeWorkingCopy(
	env *env.Env,
	remoteArg string,
) (remote repo.WorkingCopy) {
	var err error

	switch c.RemoteType {
	case repo.RemoteTypeNativeDotenvXDG:
		if remote, err = local_working_copy.MakeFromConfigAndXDGDotenvPath(
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
