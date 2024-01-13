package files

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func Rename(src, dst string) (err error) {
	if err = os.Rename(src, dst); err != nil {
		err = errors.Wrapf(err, "Src: %q, Dst: %q", src, dst)
		return
	}

	return
}
