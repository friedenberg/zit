package repo_remote

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"syscall"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/delim_io"
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

	roundTripper.Args = []string{
		"zit",
		"serve",
	}

	if remote.GetConfig().Verbose {
		roundTripper.Args = append(roundTripper.Args, "-verbose")
	}

	roundTripper.Args = append(roundTripper.Args, "-")

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

	roundTripper.Args = []string{
		"ssh",
		arg,
		"zit",
		"serve",
	}

	if remote.GetConfig().Verbose {
		roundTripper.Args = append(roundTripper.Args, "-verbose")
	}

	roundTripper.Args = append(roundTripper.Args, "-")

	if err = roundTripper.initialize(remote); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (roundTripper *HTTPRoundTripperStdio) initialize(
	remote *repo_local.Repo,
) (err error) {
	// roundTripper.Stderr = os.Stderr
	var stderrReadCloser io.ReadCloser

	if stderrReadCloser, err = roundTripper.StderrPipe(); err != nil {
		err = errors.Wrap(err)
		return
	}

	go func() {
		if _, err = delim_io.CopyWithPrefixOnDelim(
			'\n',
			"remote",
      remote.GetUI(),
			stderrReadCloser,
			false,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}()

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

	remote.After(roundTripper.cancel)

	return
}

func (roundTripper *HTTPRoundTripperStdio) cancel() (err error) {
	if roundTripper.Process != nil {
		if err = roundTripper.WriteCloser.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = roundTripper.Process.Signal(syscall.SIGHUP); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = roundTripper.Wait(); err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
