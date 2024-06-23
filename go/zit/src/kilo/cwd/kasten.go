package cwd

import (
	"os"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type Kasten = sku.ExternalMaybe

func (c *CwdFiles) tryKasten(fi os.FileInfo, dir string) (err error) {
	var h kennung.Kasten
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

	t, ok := c.kisten.Get(h.String())

	if !ok {
		t = &sku.ExternalMaybe{}
	}

	if err = t.Kennung.SetWithKennung(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.FDs.Objekte.ResetWith(f)

	return c.kisten.Add(t)
}
