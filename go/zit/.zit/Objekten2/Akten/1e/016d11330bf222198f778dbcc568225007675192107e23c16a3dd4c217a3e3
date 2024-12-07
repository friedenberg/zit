package files

import (
	"os"
	"os/exec"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"golang.org/x/term"
)

func IsTty(f *os.File) (ok bool) {
	ok = term.IsTerminal(int(f.Fd()))
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
