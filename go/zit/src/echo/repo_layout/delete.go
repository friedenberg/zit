package repo_layout

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

func (s Layout) DeleteAll(p string) (err error) {
	if s.IsDryRun() {
		return
	}

	if err = os.RemoveAll(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Layout) Delete(p string) (err error) {
	if s.IsDryRun() {
		return
	}

	if err = os.Remove(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
