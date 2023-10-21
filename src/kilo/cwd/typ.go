package cwd

import (
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type Typ = sku.ExternalMaybe

func (c *CwdFiles) tryTyp(fi os.FileInfo, dir string) (err error) {
	var h kennung.Typ
	var f fd.FD

	if f, err = fd.FileInfo(fi, dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	pathMinusExt := strings.ToLower(f.FileNameSansExt())

	if err = h.Set(pathMinusExt); err != nil {
		err = errors.Wrap(err)
		return
	}

	t, _ := c.Typen.Get(h.String())

	if err = t.Kennung.SetWithKennung(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.FDs.Objekte = f

	return c.Typen.Add(t)
}
