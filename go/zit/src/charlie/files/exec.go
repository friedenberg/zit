package files

import (
	"os"
	"os/exec"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"golang.org/x/sys/unix"
)

func IsTty(f *os.File) (ok bool) {
	_, err := unix.IoctlGetTermios(int(f.Fd()), unix.TIOCGETA)
	ok = err == nil

	return
}

func OpenFiles(p ...string) (err error) {
	cmd := exec.Command("open", p...)

	if err = cmd.Run(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
