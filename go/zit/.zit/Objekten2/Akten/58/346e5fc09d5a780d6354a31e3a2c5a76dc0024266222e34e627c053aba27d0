package files

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

func Rename(src, dst string) (err error) {
	if err = os.Rename(src, dst); err != nil {
		err = errors.Wrapf(err, "Src: %q, Dst: %q", src, dst)
		return
	}

	return
}
