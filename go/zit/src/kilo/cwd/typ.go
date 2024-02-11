package cwd

import (
	"os"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Typ = sku.ExternalMaybe

func (c *CwdFiles) tryTyp(fi os.FileInfo, dir string) (err error) {
	var h kennung.Typ
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

	t, ok := c.Typen.Get(h.String())

	if !ok {
		t = &sku.ExternalMaybe{}
	}

	if err = t.Kennung.SetWithKennung(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.FDs.Objekte.ResetWith(f)

	return c.Typen.Add(t)
}
