package open_file_guard

import (
	"os"
	"os/exec"
)

func IsTty() (ok bool) {
	cmd := exec.Command("tty", "-s")
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err == nil {
		ok = true
	}

	return
}

func OpenVimWithArgs(args []string, files ...string) (err error) {
	var cmd *exec.Cmd

	args = append(args, "-p")

	if IsTty() {
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
		err = _Error(err)
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
		err = _Error(err)
		return
	}

	return
}
