package exec_editor

import (
	"os"
	"os/exec"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/files"
)

func OpenVimWithArgs(args []string, fs ...string) (err error) {
	var cmd *exec.Cmd

	if len(fs) == 0 {
		err = errors.Wrap(files.ErrEmptyFileList)
		return
	}

	args = append(args, "-p")

	if files.IsTty(os.Stdin) {
		cmd = exec.Command(
			GetEditor(),
			append(args, fs...)...,
		)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd = exec.Command(
			"mvim",
			append(append(args, "-f"), fs...)...,
		)
	}

	if err = cmd.Run(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func GetEditor() string {
	var ed string

	if ed = os.Getenv("EDITOR"); ed != "" {
		return ed
	}

	if ed = os.Getenv("VISUAL"); ed != "" {
		return ed
	}

	return "vim"
}

func EditorIsVim() bool {
	switch GetEditor() {
	case "vim", "nvim":
		return true

	default:
		return false
	}
}

func OpenEditor(p ...string) (err error) {
	if EditorIsVim() {
		return OpenVimWithArgs(
			nil,
			p...,
		)
	} else {
		panic("not implemented")
	}
}
