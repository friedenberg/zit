package script_config

import (
	"fmt"
	"os/exec"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type ScriptConfig struct {
	Description string            `toml:"description"`
	Shell       []string          `toml:"shell,omitempty"`
	Script      string            `toml:"script,omitempty,multiline"`
	Env         map[string]string `toml:"env,omitempty"`
}

func (a *ScriptConfig) Environ() map[string]string {
	return a.Env
}

func (a *ScriptConfig) Equals(b *ScriptConfig) bool {
	if len(a.Shell) != len(b.Shell) {
		return false
	}

	if a.Description != b.Description {
		return false
	}

	for k, v := range a.Shell {
		v1 := b.Shell[k]

		if v != v1 {
			return false
		}
	}

	if a.Script != b.Script {
		return false
	}

	if len(a.Env) != len(b.Env) {
		return false
	}

	for k, v := range a.Env {
		v1, ok := b.Env[k]

		if !ok {
			return false
		}

		if v != v1 {
			return false
		}
	}

	return true
}

func (s *ScriptConfig) Merge(s2 *ScriptConfig) {
	if s2 == nil {
		return
	}

	if s2.Description != "" {
		s.Description = s2.Description
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
		all := append(s.Shell, s.Script, "--")
		all = append(all, args...)
		c = exec.Command(s.Shell[0], all[1:]...)

	case s.Script != "":
		all := []string{
			"--noprofile",
			"--norc",
			"-c",
		}

		all = append(all, s.Script, "--")
		all = append(all, args...)
		c = exec.Command("bash", all...)

	case len(s.Shell) > 0:
		all := append(s.Shell, "--")
		all = append(all, args...)

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
