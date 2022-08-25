package store_checkout

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/age_io"
)

func (s Store) AkteShaFromPath(p string) (sh sha.Sha, err error) {
	var aw age_io.Writer

	if aw, err = s.storeZettel.AkteWriter(); err != nil {
		err = errors.Error(err)
		return
	}

	var f *os.File

	if f, err = open_file_guard.Open(p); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(f)

	if _, err = io.Copy(aw, f); err != nil {
		err = errors.Error(err)
		return
	}

	sh = aw.Sha()

	return
}
