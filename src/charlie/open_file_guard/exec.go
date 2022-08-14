package open_file_guard

import (
	"os"
	"os/exec"

	"github.com/friedenberg/zit/src/bravo/errors"
	"golang.org/x/sys/unix"
)

func IsTty(f *os.File) (ok bool) {
	_, err := unix.IoctlGetTermios(int(f.Fd()), unix.TIOCGETA)
	ok = err == nil

	return
}

func OpenVimWithArgs(args []string, files ...string) (err error) {
	var cmd *exec.Cmd

	args = append(args, "-p")

	if IsTty(os.Stdin) {
		cmd = exec.Command(
			"vim",
			append(args, files...)...,
		)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd = exec.Command(
			"mvim",
			append(append(args, "-f"), files...)...,
		)
	}

	if err = cmd.Run(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func OpenEditor(p ...string) (err error) {
	return OpenVimWithArgs(
		nil,
		p...,
	)
}

func OpenFiles(p ...string) (err error) {
	cmd := exec.Command("open", p...)

	if err = cmd.Run(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
