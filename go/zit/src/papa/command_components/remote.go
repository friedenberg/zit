package command_components

import (
	"flag"
	"fmt"
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
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

	return cmd.MakeRemoteWorkingCopy(env, f.Args()[0])
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
		if remote, err = cmd.MakeFromConfigAndXDGDotenvPath(
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
		if remote, err = cmd.MakeRemoteHTTPFromXDGDotenvPath(
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

func (cmd Remote) MakeRemoteWorkingCopy(
	env *env.Env,
	remoteArg string,
) (remote repo.WorkingCopy) {
	var err error

	switch cmd.RemoteType {
	case repo.RemoteTypeNativeDotenvXDG:
		if remote, err = cmd.MakeFromConfigAndXDGDotenvPath(
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
		if remote, err = cmd.MakeRemoteHTTPFromXDGDotenvPath(
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

func (cmd *Remote) MakeRemoteHTTPFromXDGDotenvPath(
	context *errors.Context,
	config config_mutable_cli.Config,
	xdgDotenvPath string,
	options env.Options,
) (remoteHTTP *repo_remote.HTTP, err error) {
	var remote *local_working_copy.Repo

	if remote, err = cmd.MakeFromConfigAndXDGDotenvPath(
		context,
		config,
		xdgDotenvPath,
		options,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	remoteHTTP = &repo_remote.HTTP{
		Repo: remote,
	}

	var httpRoundTripper repo_remote.HTTPRoundTripperUnixSocket

	if err = httpRoundTripper.Initialize(remote); err != nil {
		err = errors.Wrap(err)
		return
	}

	remoteHTTP.Client = http.Client{
		Transport: &httpRoundTripper,
	}

	go func() {
		if err := remote.Serve(httpRoundTripper.UnixSocket); err != nil {
			remote.CancelWithError(err)
		}
	}()

	return
}
