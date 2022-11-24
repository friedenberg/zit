package cwd_files

import (
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/typ"
	"github.com/friedenberg/zit/src/fd"
)

type CwdTyp = typ.External

func (c *CwdFiles) tryTyp(fi os.FileInfo) (err error) {
	var h kennung.Typ

	fd := fd.FileInfo(fi)
	pathMinusExt := strings.ToLower(fd.FileNameSansExt())

	if err = h.Set(pathMinusExt); err != nil {
		err = errors.Wrap(err)
		return
	}

	var t *CwdTyp

	ok := false

	if t, ok = c.Typen[pathMinusExt]; !ok {
		t = &CwdTyp{
			Named: typ.Named{
				Kennung: h,
			},
		}
	}

	t.FD = fd
	c.Typen[h.String()] = t

	return
}
