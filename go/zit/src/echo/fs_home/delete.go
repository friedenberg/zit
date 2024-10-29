package fs_home

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

func (s Home) DeleteAll(p string) (err error) {
	if s.dryRun {
		return
	}

	if err = os.RemoveAll(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Home) Delete(p string) (err error) {
	if s.dryRun {
		return
	}

	if err = os.Remove(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
