package cwd

import (
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type Zettel = sku.ExternalMaybe

func (c *CwdFiles) tryZettel(d string, a string, p string) (err error) {
	var fd kennung.FD

	if fd, err = kennung.FDFromPath(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	var h kennung.Hinweis

	if h, err = fd.GetHinweis(); err != nil {
		err = errors.Wrap(err)
		return
	}

	t, _ := c.Zettelen.Get(h.String())
	t.Kennung = kennung.Kennung2{KennungPtr: &h}

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
