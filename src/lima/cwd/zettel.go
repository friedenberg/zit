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
	t.Kennung = &h

	if path.Ext(a) == c.erworben.GetZettelFileExtension() {
		t.FDs.Objekte.Path = p
	} else {
		t.FDs.Akte.Path = p
	}

	return c.Zettelen.Add(t)
}
