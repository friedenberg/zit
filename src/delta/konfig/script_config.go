package konfig

import (
	"fmt"
	"os/exec"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type ScriptConfig struct {
	Shell  []string
	Script string
	Env    map[string]string
}

func (s *ScriptConfig) Merge(s2 *ScriptConfig) {
	if s2 == nil {
		return
	}

	if len(s2.Shell) > 0 {
		s.Shell = s2.Shell
	}

	if s2.Script != "" {
		s.Script = s2.Script
	}

	if len(s.Env) == 0 {
		s.Env = make(map[string]string)
	}

	for k, v := range s2.Env {
		s.Env[k] = v
	}
}

func (s ScriptConfig) Cmd(args ...string) (c *exec.Cmd, err error) {
	switch {
	case s.Script == "" && len(s.Shell) == 0:
		err = errors.Errorf("no script or shell set")
		return

	case s.Script != "" && len(s.Shell) > 0:
		all := append(s.Shell, s.Script)
		c = exec.Command(s.Shell[0], all[1:]...)

	case s.Script != "":
		all := []string{
			"--noprofile",
			"--norc",
			"-c",
		}

		all = append(all, s.Script)
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

	envCollapsed := make([]string, 0, len(s.Env))

	for k, v := range s.Env {
		envCollapsed = append(envCollapsed, fmt.Sprintf("%s=%s", k, v))
	}

	c.Env = envCollapsed

	return
}
