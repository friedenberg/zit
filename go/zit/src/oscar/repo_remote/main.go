package repo_remote

import (
	"net/http"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

func MakeRemoteHTTPFromXDGDotenvPath(
	context *errors.Context,
	config config_mutable_cli.Config,
	xdgDotenvPath string,
	options env.Options,
) (remoteHTTP *HTTP, err error) {
	var remote *local_working_copy.Repo

	if remote, err = local_working_copy.MakeFromConfigAndXDGDotenvPath(
		context,
		config,
		xdgDotenvPath,
		options,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	remoteHTTP = &HTTP{
		Repo: remote,
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
	remote := local_working_copy.Make(
		env,
		local_working_copy.OptionsEmpty,
	)

	remoteHTTP = &HTTP{
		Repo: remote,
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
	remote := local_working_copy.Make(
		env,
		local_working_copy.OptionsEmpty,
	)

	remoteHTTP = &HTTP{
		Repo: remote,
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
