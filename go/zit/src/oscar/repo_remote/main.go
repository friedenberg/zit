package repo_remote

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

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
