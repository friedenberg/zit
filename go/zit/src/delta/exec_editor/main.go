package exec_editor

import (
	"os"
	"os/exec"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
)

type Editor struct {
	path string
	name string
	tipe Type

	options []string

	ui interfaces.FuncIter[string]
}

func getEditor() string {
	var ed string

	if ed = os.Getenv("EDITOR"); ed != "" {
		return ed
	}

	if ed = os.Getenv("VISUAL"); ed != "" {
		return ed
	}

	return "vim"
}

func MakeEditorWithVimOptions(
	ph interfaces.FuncIter[string],
	options []string,
) Editor {
	return MakeEditor(
		ph,
		map[Type][]string{
			TypeVim: options,
		},
	)
}

func MakeEditor(
	ph interfaces.FuncIter[string],
	options map[Type][]string,
) Editor {
	editor := Editor{
		path: getEditor(),
	}

	editor.name = filepath.Base(editor.path)

	switch editor.name {
	case "vim", "nvim":
		editor.tipe = TypeVim
	}

	editor.options = options[editor.tipe]

	return editor
}

func (c Editor) Run(
	files []string,
) (err error) {
	switch c.tipe {
	case TypeVim:
		if err = c.runVim(files); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		if err = c.runUnknown(files); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c Editor) runUnknown(
	files []string,
) (err error) {
	return
}

func (c Editor) runVim(
	files []string,
) (err error) {
	vimOptions := c.options
	vimArgs := make([]string, 0, (len(vimOptions)*2)+1)
	vimArgs = append(vimArgs, "-f")

	for _, o := range vimOptions {
		vimArgs = append(vimArgs, "-c", o)
	}

	v := "vim started"

	if err = c.ui(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.openWithArgs(vimArgs, files...); err != nil {
		err = errors.Wrap(err)
		return
	}

	v = "vim exited"

	if err = c.ui(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (editor Editor) openWithArgs(args []string, fs ...string) (err error) {
	if len(fs) == 0 {
		err = errors.Wrap(files.ErrEmptyFileList)
		return
	}

	allArgs := append(args, fs...)

	cmd := exec.Command(
		editor.path,
		allArgs...,
	)

	if files.IsTty(os.Stdin) {
		cmd.Stdin = os.Stdin
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		err = errors.Wrapf(err, "Cmd: %s", cmd)
		return
	}

	return
}
