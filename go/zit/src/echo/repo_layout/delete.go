package repo_layout

import (
	"os"
	"path/filepath"

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
  p = filepath.Clean(p)

  if p == "." {
    err = errors.Errorf("invalid delete request: %q", p)
    return
  }

	if s.IsDryRun() {
		return
	}

	if err = os.Remove(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
