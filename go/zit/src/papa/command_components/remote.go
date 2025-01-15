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
	LocalWorkingCopy
	LocalArchive
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

// TODO
func (cmd Remote) MakeArchiveFromFlagSet(
	env *env.Env,
	f *flag.FlagSet,
) (remote repo.Archive) {
	remoteArg := f.Args()[0]

	var err error

	switch cmd.RemoteType {
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
				cmd.RemoteType,
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
				cmd.RemoteType,
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
				cmd.RemoteType,
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
				cmd.RemoteType,
				remoteArg,
			)
		}

	default:
		env.CancelWithNotImplemented()
	}

	return
}

func (cmd Remote) MakeWorkingCopy(
	env *env.Env,
	remoteArg string,
) (remote repo.WorkingCopy) {
	var err error

	switch cmd.RemoteType {
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
				cmd.RemoteType,
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
				cmd.RemoteType,
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
				cmd.RemoteType,
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
				cmd.RemoteType,
				remoteArg,
			)
		}

	default:
		env.CancelWithNotImplemented()
	}

	return
}
