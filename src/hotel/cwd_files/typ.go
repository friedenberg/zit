package cwd_files

import (
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/fd"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
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
			Sku: sku.External[kennung.Typ, *kennung.Typ]{
				Kennung: h,
			},
		}
	}

	t.FD = fd
	c.Typen[h.String()] = t

	return
}
