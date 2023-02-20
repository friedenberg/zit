package cwd

import (
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/hotel/etikett"
)

func (c *CwdFiles) tryEtikett(fi os.FileInfo) (err error) {
	var h kennung.Etikett

	fd := kennung.FileInfo(fi)
	pathMinusExt := strings.ToLower(fd.FileNameSansExt())

	if err = h.Set(pathMinusExt); err != nil {
		err = errors.Wrap(err)
		return
	}

	var t *etikett.External

	ok := false

	if t, ok = c.Etiketten[h]; !ok {
		t = &etikett.External{
			Sku: sku.External[kennung.Etikett, *kennung.Etikett]{
				Kennung: h,
			},
		}
	}

	t.Sku.ObjekteFD = fd
	c.Etiketten[h] = t

	return
}
