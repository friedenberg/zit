package konfig

import (
	"os/exec"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type ScriptConfig struct {
	Shell  []string
	Script string
}

func (s ScriptConfig) Cmd(args ...string) (c *exec.Cmd, err error) {
	switch {
	case s.Script == "" && len(s.Shell) == 0:
		err = errors.Errorf("no script or shell set")
		return

	case s.Script != "":
		all := []string{
			"--noprofile",
			"--norc",
			"-c",
		}

		all = append(all, args...)
		c = exec.Command("bash", all...)

	case len(s.Shell) > 0:
		all := append(s.Shell, args...)
		if len(all) > 1 {
			c = exec.Command(all[0], all[1:]...)
		} else {
			c = exec.Command(all[0])
		}
	}

	return
}