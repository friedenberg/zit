package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/exec_editor"
)

type Open struct {
	Options []string
}

func (c Open) Run(
	ph interfaces.FuncIter[string],
	args ...string,
) (err error) {
	var editor exec_editor.Editor

	if editor, err = exec_editor.MakeEditorWithVimOptions(
		ph,
		c.Options,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = editor.Run(args); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
