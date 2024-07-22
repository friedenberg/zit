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
	vimArgs := make([]string, 0, len(c.Options)*2)

	for _, o := range c.Options {
		vimArgs = append(vimArgs, "-c", o)
	}

	v := "vim started"

	if err = ph(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = exec_editor.OpenVimWithArgs(vimArgs, args...); err != nil {
		err = errors.Wrap(err)
		return
	}

	v = "vim exited"

	if err = ph(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
