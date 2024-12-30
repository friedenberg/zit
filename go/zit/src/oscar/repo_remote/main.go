package repo_remote

import (
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

func MakeRemoteHTTPFromXDGDotenvPath(
	context *errors.Context,
	config config_mutable_cli.Config,
	xdgDotenvPath string,
	options env.Options,
) (remoteHTTP *HTTP, err error) {
	var remote *repo_local.Repo

	if remote, err = repo_local.MakeFromConfigAndXDGDotenvPath(
		context,
		config,
		xdgDotenvPath,
		options,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	remoteHTTP = &HTTP{
		remote: remote,
	}

	var httpRoundTripper HTTPRoundTripperUnixSocket

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

func MakeRemoteStdioLocal(
	env *env.Env,
	dir string,
) (remoteHTTP *HTTP, err error) {
	remote := repo_local.Make(
		env,
		repo_local.OptionsEmpty,
	)

	remoteHTTP = &HTTP{
		remote: remote,
	}

	var httpRoundTripper HTTPRoundTripperStdio

	httpRoundTripper.Dir = dir

	if err = httpRoundTripper.InitializeWithLocal(remote); err != nil {
		err = errors.Wrap(err)
		return
	}

	remoteHTTP.Client.Transport = &httpRoundTripper

	return
}

func MakeRemoteStdioSSH(
	env *env.Env,
	arg string,
) (remoteHTTP *HTTP, err error) {
	remote := repo_local.Make(
		env,
		repo_local.OptionsEmpty,
	)

	remoteHTTP = &HTTP{
		remote: remote,
	}

	var httpRoundTripper HTTPRoundTripperStdio

	if err = httpRoundTripper.InitializeWithSSH(
		remote,
		arg,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	remoteHTTP.Client.Transport = &httpRoundTripper

	return
}
