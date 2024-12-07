package dir_layout

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

func (s DirLayout) DeleteAll(p string) (err error) {
	if s.dryRun {
		return
	}

	if err = os.RemoveAll(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s DirLayout) Delete(p string) (err error) {
	if s.dryRun {
		return
	}

	if err = os.Remove(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
