package konfig

import (
	"os/exec"
)

type RemoteScript interface {
	Cmd(args ...string) (*exec.Cmd, error)
}
