package cwd

import (
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type Etikett = sku.ExternalMaybe

func (c *CwdFiles) tryEtikett(fi os.FileInfo, dir string) (err error) {
	var h kennung.Etikett
	var f *fd.FD

	if f, err = fd.FileInfo(fi, dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	pathMinusExt := strings.ToLower(f.FileNameSansExt())

	if err = h.Set(pathMinusExt); err != nil {
		err = errors.Wrap(err)
		return
	}

	t, ok := c.Etiketten.Get(h.String())

	if !ok {
		t = &sku.ExternalMaybe{}
	}

	if err = t.Kennung.SetWithKennung(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.FDs.Objekte.ResetWith(f)

	return c.Etiketten.Add(t)
}