package remote_http

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"syscall"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/delim_io"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
)

type RoundTripperStdio struct {
	exec.Cmd
	io.WriteCloser
	io.ReadCloser
	roundTripperBufio
}

func (roundTripper *RoundTripperStdio) InitializeWithLocal(
	envUI env_ui.Env,
) (err error) {
	if roundTripper.Path, err = exec.LookPath("zit"); err != nil {
		err = errors.Wrap(err)
		return
	}

	roundTripper.Args = []string{
		"zit",
		"serve",
		"-print-time=false", // TODO switch to passing this from envUI.GetCLIConfig()
	}

	if envUI.GetCLIConfig().Verbose {
		roundTripper.Args = append(roundTripper.Args, "-verbose")
	}

	roundTripper.Args = append(roundTripper.Args, "-")

	if err = roundTripper.initialize(envUI); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (roundTripper *RoundTripperStdio) InitializeWithSSH(
	envUI env_ui.Env,
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

	if envUI.GetCLIConfig().Verbose {
		roundTripper.Args = append(roundTripper.Args, "-verbose")
	}

	roundTripper.Args = append(roundTripper.Args, "-")

	if err = roundTripper.initialize(envUI); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (roundTripper *RoundTripperStdio) initialize(
	envUI env_ui.Env,
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
			envUI.GetUI(),
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

	envUI.After(roundTripper.cancel)

	return
}

func (roundTripper *RoundTripperStdio) cancel() (err error) {
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
