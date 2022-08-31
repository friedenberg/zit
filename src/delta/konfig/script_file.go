package konfig

import "os/exec"

type RemoteScriptFile struct {
	Path string
}

func (s RemoteScriptFile) Cmd(args ...string) (c *exec.Cmd, err error) {
	c = exec.Command(s.Path, args...)

	return
}
