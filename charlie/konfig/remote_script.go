package konfig

import "os/exec"

type RemoteScript interface {
	Cmd(args []string) (*exec.Cmd, error)
}

type RemoteScriptConfig struct {
	SupportedTypes    []_Type
	SupportedCommands []string
	Shell             string
	Script            string
}

type RemoteScriptFile struct {
	Path string
}

func (s RemoteScriptFile) Cmd(args []string) (c *exec.Cmd, err error) {
	c = exec.Command(s.Path, args...)

	return
}

func (s RemoteScriptConfig) Cmd(args []string) (c *exec.Cmd, err error) {
	shell := s.Shell

	if shell == "" {
		shell = "/bin/bash"
	}

	c = exec.Command(
		shell,
		append(
			[]string{
				"-c",
				s.Script,
			},
			args...,
		)...,
	)

	return
}
