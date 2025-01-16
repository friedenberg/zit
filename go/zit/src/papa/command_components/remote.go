package command_components

import (
	"flag"
	"fmt"
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/oscar/repo_remote"
)

type Remote struct {
	Env
	RepoLayout
	LocalWorkingCopy
	LocalArchive

	RemoteType repo.RemoteType
}

func (cmd *Remote) SetFlagSet(f *flag.FlagSet) {
	f.Var(&cmd.RemoteType, "remote-type", fmt.Sprintf("%s", repo.GetAllRemoteTypes()))
}

func (cmd Remote) MakeWorkingCopyFromFlagSet(
	req command.Request,
) (remote repo.WorkingCopy) {
	if len(req.Args()) == 0 {
		// TODO add info about remote options
		req.CancelWithBadRequestf("requires a remote to be specified")
	}

	return cmd.MakeRemoteWorkingCopy(req, req.Args()[0])
}

// TODO
func (cmd Remote) MakeArchiveFromFlagSet(
	req command.Request,
) (remote repo.Archive) {
	remoteArg := req.Args()[0]
	env := cmd.MakeEnv(req)

	var err error

	switch cmd.RemoteType {
	case repo.RemoteTypeNativeDotenvXDG:
		remote = cmd.LocalWorkingCopy.MakeLocalWorkingCopyFromConfigAndXDGDotenvPath(
			req,
			remoteArg,
			env.GetOptions(),
		)

	case repo.RemoteTypeStdioLocal:
		remote = repo_remote.MakeRemoteStdioLocal(
			env,
			remoteArg,
		)

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
			req,
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
	req command.Request,
	remoteArg string,
) (remote repo.WorkingCopy) {
	var err error

	switch cmd.RemoteType {
	case repo.RemoteTypeNativeDotenvXDG:
		remote = cmd.LocalWorkingCopy.MakeLocalWorkingCopyFromConfigAndXDGDotenvPath(
			req,
			remoteArg,
			env.Options{},
		)

	case repo.RemoteTypeStdioLocal:
		remote = repo_remote.MakeRemoteStdioLocal(
			cmd.MakeEnv(req),
			remoteArg,
		)

	case repo.RemoteTypeStdioSSH:
		if remote, err = repo_remote.MakeRemoteStdioSSH(
			cmd.MakeEnv(req),
			remoteArg,
		); err != nil {
			req.CancelWithErrorAndFormat(
				err,
				"RemoteType: %q, Remote: %q",
				cmd.RemoteType,
				remoteArg,
			)
		}

	case repo.RemoteTypeSocketUnix:
		if remote, err = cmd.MakeRemoteHTTPFromXDGDotenvPath(
			req,
			remoteArg,
			env.Options{},
		); err != nil {
			req.CancelWithErrorAndFormat(
				err,
				"RemoteType: %q, Remote: %q",
				cmd.RemoteType,
				remoteArg,
			)
		}

	default:
		req.CancelWithNotImplemented()
	}

	return
}

func (cmd *Remote) MakeLocalWorkingCopyFromConfigAndXDGDotenvPath(
	req command.Request,
	xdgDotenvPath string,
	options env.Options,
) repo.Archive {
	repoLayout := cmd.MakeRepoLayout(req, false)

	return cmd.MakeLocalArchive(repoLayout)
}

func (cmd *Remote) MakeRemoteHTTPFromXDGDotenvPath(
	req command.Request,
	xdgDotenvPath string,
	options env.Options,
) (remoteHTTP *repo_remote.HTTP, err error) {
	remote := cmd.LocalWorkingCopy.MakeLocalWorkingCopyFromConfigAndXDGDotenvPath(
		req,
		xdgDotenvPath,
		options,
	)

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
			req.CancelWithError(err)
		}
	}()

	return
}
