package cwd

import (
	"path"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type Zettel = sku.KennungFDPair

func (c *CwdFiles) tryZettel(
	dir string,
	name string,
	fullPath string,
	unsure bool,
) (err error) {
	var f *fd.FD

	if f, err = fd.FDFromPath(fullPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	var h kennung.Hinweis

	if h, err = kennung.GetHinweis(f, unsure); err != nil {
		err = errors.Wrap(err)
		return
	}

	t, ok := c.zettelen.Get(h.String())

	if !ok {
		t = &sku.KennungFDPair{}
	}

	if err = t.Kennung.SetWithKennung(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	if path.Ext(name) == c.erworben.GetZettelFileExtension() {
		if err = t.FDs.Objekte.SetPath(fullPath); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = t.FDs.Akte.SetPath(fullPath); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if unsure {
		if err = c.unsureZettelen.Add(t); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = c.zettelen.Add(t); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
