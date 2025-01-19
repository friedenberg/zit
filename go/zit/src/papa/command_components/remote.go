package command_components

import (
	"flag"
	"fmt"
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
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
			req,
			env,
			remoteArg,
		)

	case repo.RemoteTypeStdioSSH:
		remote = cmd.MakeRemoteStdioSSH(
			req,
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
			env_ui.Options{},
		)

	case repo.RemoteTypeStdioLocal:
		remote = cmd.MakeRemoteStdioLocal(
			req,
			cmd.MakeEnv(req),
			remoteArg,
		)

	case repo.RemoteTypeStdioSSH:
		remote = cmd.MakeRemoteStdioSSH(
			req,
			cmd.MakeEnv(req),
			remoteArg,
		)

	case repo.RemoteTypeSocketUnix:
		remote = cmd.MakeRemoteHTTPFromXDGDotenvPath(
			req,
			remoteArg,
			env_ui.Options{},
		)

	default:
		req.CancelWithNotImplemented()
	}

	return
}

func (cmd *Remote) MakeRemoteHTTPFromXDGDotenvPath(
	req command.Request,
	xdgDotenvPath string,
	options env_ui.Options,
) (remoteHTTP *remote_http.Client) {
	envLocal := cmd.MakeEnvWithXDGLayoutAndOptions(
		req,
		xdgDotenvPath,
		options,
	)

	envRepo := cmd.MakeRepoLayoutFromEnvLocal(envLocal)

	remote := cmd.MakeLocalArchive(envRepo)

	server := remote_http.Server{
		EnvLocal: envLocal,
		Repo:     remote,
	}

	remoteHTTP = &remote_http.Client{
		Repo: remote,
	}

	var httpRoundTripper remote_http.RoundTripperUnixSocket

	if err := httpRoundTripper.Initialize(server); err != nil {
		req.CancelWithError(err)
	}

	remoteHTTP.Client = http.Client{
		Transport: &httpRoundTripper,
	}

	go func() {
		if err := server.Serve(httpRoundTripper.UnixSocket); err != nil {
			req.CancelWithError(err)
		}
	}()

	return
}

func (cmd *Remote) MakeRemoteStdioSSH(
	req command.Request,
	env env_local.Env,
	arg string,
) (remoteHTTP *remote_http.Client) {
	repoLayout := cmd.MakeRepoLayout(req, false)
	remote := cmd.MakeLocalArchive(repoLayout)

	remoteHTTP = &remote_http.Client{
		Repo: remote,
	}

	var httpRoundTripper remote_http.RoundTripperStdio

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
	req command.Request,
	env env_local.Env,
	dir string,
) (remoteHTTP *remote_http.Client) {
	repoLayout := cmd.MakeRepoLayout(req, false)
	remote := cmd.MakeLocalArchive(repoLayout)

	remoteHTTP = &remote_http.Client{
		Repo: remote,
	}

	var httpRoundTripper remote_http.RoundTripperStdio

	httpRoundTripper.Dir = dir

	if err := httpRoundTripper.InitializeWithLocal(remote); err != nil {
		env.CancelWithError(err)
	}

	remoteHTTP.Client.Transport = &httpRoundTripper

	return
}
