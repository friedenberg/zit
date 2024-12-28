package repo_remote

import (
	"bufio"
	"io"
	"os"
	"os/exec"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type HTTPRoundTripperStdio struct {
	exec.Cmd
	io.WriteCloser
	io.ReadCloser
	HTTPRoundTripperBufio
}

func (roundTripper *HTTPRoundTripperStdio) InitializeWithLocal(
	remote *repo_local.Repo,
) (err error) {
	if roundTripper.Path, err = exec.LookPath("zit"); err != nil {
		err = errors.Wrap(err)
		return
	}

	roundTripper.Stderr = os.Stderr

	roundTripper.Args = []string{
		"zit",
		"serve",
		// "-verbose",
		"-",
	}

	if err = roundTripper.initialize(remote); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (roundTripper *HTTPRoundTripperStdio) InitializeWithSSH(
	remote *repo_local.Repo,
	arg string,
) (err error) {
	if roundTripper.Path, err = exec.LookPath("ssh"); err != nil {
		err = errors.Wrap(err)
		return
	}

	roundTripper.Stderr = os.Stderr

	roundTripper.Args = []string{
		arg,
		"zit",
		"serve",
		// "-verbose",
		"-",
	}

	if err = roundTripper.initialize(remote); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (roundTripper *HTTPRoundTripperStdio) initialize(
	remote *repo_local.Repo,
) (err error) {
	if roundTripper.WriteCloser, err = roundTripper.StdinPipe(); err != nil {
		err = errors.Wrap(err)
		return
	}

	roundTripper.Writer = bufio.NewWriter(roundTripper.WriteCloser)

	if roundTripper.ReadCloser, err = roundTripper.StdoutPipe(); err != nil {
		err = errors.Wrap(err)
		return
	}

	roundTripper.Reader = bufio.NewReader(roundTripper.ReadCloser)

	if err = roundTripper.Start(); err != nil {
		err = errors.Wrapf(err, "%#v", roundTripper.Cmd)
		return
	}

	return
}
