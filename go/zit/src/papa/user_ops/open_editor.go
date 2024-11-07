package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/exec_editor"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type OpenEditor struct {
	VimOptions []string
}

func (c OpenEditor) Run(
	u *env.Env,
	args ...string,
) (err error) {
	editor := exec_editor.MakeEditorWithVimOptions(
		u.PrinterHeader(),
		c.VimOptions,
	)

	if err = editor.Run(args); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
