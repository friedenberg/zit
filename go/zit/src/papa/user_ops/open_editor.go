package user_ops

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/editor"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type OpenEditor struct {
	VimOptions []string
}

func (c OpenEditor) Run(
	u *env.Local,
	args ...string,
) (err error) {
	var e editor.Editor

	if e, err = editor.MakeEditorWithVimOptions(
		u.PrinterHeader(),
		c.VimOptions,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.Run(args); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
