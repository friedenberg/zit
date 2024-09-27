package script_config

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

type RemoteScript interface {
	Cmd(args ...string) (*exec.Cmd, error)
}

type RemoteScriptWithEnv interface {
	RemoteScript
	Environ() map[string]string
}

type writerTo struct {
	cmd *exec.Cmd
}

func MakeWriterTo(
	rs RemoteScript,
	env map[string]string,
	args ...string,
) (wt *writerTo, err error) {
	wt = &writerTo{}

	if rs == nil {
		err = errors.Errorf("empty remote script")
		return
	}

	if wt.cmd, err = rs.Cmd(args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	ui.Log().Print(wt.cmd)

	envCollapsed := os.Environ()

	for k, v := range env {
		envCollapsed = append(envCollapsed, fmt.Sprintf("%s=%s", k, v))
	}

	if rswe, ok := rs.(RemoteScriptWithEnv); ok {
		for k, v := range rswe.Environ() {
			envCollapsed = append(envCollapsed, fmt.Sprintf("%s=%s", k, v))
		}
	}

	ui.TodoP2("determine how stderr and env should be handled")
	wt.cmd.Stderr = os.Stderr
	wt.cmd.Env = envCollapsed

	return
}

func MakeWriterToWithStdin(
	rs RemoteScript,
	env map[string]string,
	r io.Reader,
	args ...string,
) (wt *writerTo, err error) {
	if wt, err = MakeWriterTo(rs, env, args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	wt.cmd.Stdin = r

	return
}

func (wt *writerTo) WriteTo(w io.Writer) (n int64, err error) {
	var r io.ReadCloser

	if r, err = wt.cmd.StdoutPipe(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = wt.cmd.Start(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if n, err = io.Copy(w, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = wt.cmd.Wait(); err != nil {
		err = errors.Wrapf(err, "Command: '%s'", wt.cmd.String())
		return
	}

	return
}
