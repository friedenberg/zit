package command_components

import (
	"flag"
	"fmt"
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/oscar/remote_http"
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

// TODO
func (cmd Remote) MakeArchive(
	req command.Request,
	remoteArg string,
) (remote repo.Repo) {
	env := cmd.MakeEnv(req)

	switch cmd.RemoteType {
	case repo.RemoteTypeNativeDotenvXDG:
		remote = cmd.LocalWorkingCopy.MakeLocalWorkingCopyFromConfigAndXDGDotenvPath(
			req,
			remoteArg,
			env.GetOptions(),
		)

	case repo.RemoteTypeStdioLocal:
		remote = cmd.MakeRemoteStdioLocal(
			env,
			remoteArg,
		)

	case repo.RemoteTypeStdioSSH:
		remote = cmd.MakeRemoteStdioSSH(
			env,
			remoteArg,
		)

	case repo.RemoteTypeSocketUnix:
		remote = cmd.MakeRemoteHTTPFromXDGDotenvPath(
			req,
			remoteArg,
			env.GetOptions(),
		)

	default:
		env.CancelWithNotImplemented()
	}

	return
}

func (cmd Remote) MakeRemoteWorkingCopy(
	req command.Request,
	remoteArg string,
) (remote repo.WorkingCopy) {
	switch cmd.RemoteType {
	case repo.RemoteTypeNativeDotenvXDG:
		remote = cmd.LocalWorkingCopy.MakeLocalWorkingCopyFromConfigAndXDGDotenvPath(
			req,
			remoteArg,
			env.Options{},
		)

	case repo.RemoteTypeStdioLocal:
		remote = cmd.MakeRemoteStdioLocal(
			cmd.MakeEnv(req),
			remoteArg,
		)

	case repo.RemoteTypeStdioSSH:
		remote = cmd.MakeRemoteStdioSSH(
			cmd.MakeEnv(req),
			remoteArg,
		)

	case repo.RemoteTypeSocketUnix:
		remote = cmd.MakeRemoteHTTPFromXDGDotenvPath(
			req,
			remoteArg,
			env.Options{},
		)

	default:
		req.CancelWithNotImplemented()
	}

	return
}

func (cmd *Remote) MakeLocalWorkingCopyFromConfigAndXDGDotenvPath(
	req command.Request,
	xdgDotenvPath string,
	options env.Options,
) repo.Repo {
	repoLayout := cmd.MakeRepoLayout(req, false)

	return cmd.MakeLocalArchive(repoLayout)
}

func (cmd *Remote) MakeRemoteHTTPFromXDGDotenvPath(
	req command.Request,
	xdgDotenvPath string,
	options env.Options,
) (remoteHTTP *remote_http.Remote) {
	remote := cmd.LocalWorkingCopy.MakeLocalWorkingCopyFromConfigAndXDGDotenvPath(
		req,
		xdgDotenvPath,
		options,
	)

	remoteHTTP = &remote_http.Remote{
		Repo: remote,
	}

	var httpRoundTripper remote_http.HTTPRoundTripperUnixSocket

	if err := httpRoundTripper.Initialize(remote); err != nil {
		req.CancelWithError(err)
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

func (cmd *Remote) MakeRemoteStdioSSH(
	env *env.Env,
	arg string,
) (remoteHTTP *remote_http.Remote) {
	remote := local_working_copy.Make(
		env,
		local_working_copy.OptionsEmpty,
	)

	remoteHTTP = &remote_http.Remote{
		Repo: remote,
	}

	var httpRoundTripper remote_http.HTTPRoundTripperStdio

	if err := httpRoundTripper.InitializeWithSSH(
		remote,
		arg,
	); err != nil {
		env.CancelWithError(err)
	}

	remoteHTTP.Client.Transport = &httpRoundTripper

	return
}

func (cmd *Remote) MakeRemoteStdioLocal(
	env *env.Env,
	dir string,
) (remoteHTTP *remote_http.Remote) {
	remote := local_working_copy.Make(
		env,
		local_working_copy.OptionsEmpty,
	)

	remoteHTTP = &remote_http.Remote{
		Repo: remote,
	}

	var httpRoundTripper remote_http.HTTPRoundTripperStdio

	httpRoundTripper.Dir = dir

	if err := httpRoundTripper.InitializeWithLocal(remote); err != nil {
		env.CancelWithError(err)
	}

	remoteHTTP.Client.Transport = &httpRoundTripper

	return
}
