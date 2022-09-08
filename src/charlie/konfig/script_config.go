package konfig

import "os/exec"

type ScriptConfig struct {
	Shell  []string
	Script string
}

func (s ScriptConfig) Cmd(args ...string) (c *exec.Cmd, err error) {
	shell := s.Shell

	if len(shell) == 0 {
		shell = []string{
			"bash",
			"--noprofile",
			"--norc",
			"-c",
		}
	}

	first := shell[0]

	if len(shell) > 1 {
		shell = shell[1:]
	} else {
		shell = []string{}
	}

	all := append(shell, args...)

	if s.Script != "" {
		all = append(all, s.Script)
	}

	c = exec.Command(first, all...)

	return
}
