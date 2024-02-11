package cwd

import (
	"path"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Zettel = sku.ExternalMaybe

func (c *CwdFiles) tryZettel(d string, a string, p string) (err error) {
	var f *fd.FD

	if f, err = fd.FDFromPath(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	var h kennung.Hinweis

	if h, err = kennung.GetHinweis(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	t, ok := c.Zettelen.Get(h.String())

	if !ok {
		t = &sku.ExternalMaybe{}
	}

	if err = t.Kennung.SetWithKennung(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	if path.Ext(a) == c.erworben.GetZettelFileExtension() {
		if err = t.FDs.Objekte.SetPath(p); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = t.FDs.Akte.SetPath(p); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return c.Zettelen.Add(t)
}
