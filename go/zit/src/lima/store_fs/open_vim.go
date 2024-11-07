package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/editor"
)

type Open struct {
	Options []string
}

func (c Open) Run(
	ph interfaces.FuncIter[string],
	args ...string,
) (err error) {
	var e editor.Editor

	if e, err = editor.MakeEditorWithVimOptions(
		ph,
		c.Options,
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
