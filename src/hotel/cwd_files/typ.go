package cwd_files

import (
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/fd"
	"github.com/friedenberg/zit/src/golf/typ"
)

func (c *CwdFiles) tryTyp(fi os.FileInfo) (err error) {
	var h kennung.Typ

	fd := fd.FileInfo(fi)
	pathMinusExt := strings.ToLower(fd.FileNameSansExt())

	if err = h.Set(pathMinusExt); err != nil {
		err = errors.Wrap(err)
		return
	}

	var t *typ.External

	ok := false

	if t, ok = c.Typen[pathMinusExt]; !ok {
		t = &typ.External{
			Named: typ.Named{
				Kennung: h,
			},
		}
	}

	t.FD = fd
	c.Typen[h.String()] = t

	return
}