package cwd

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type Zettel = sku.ExternalMaybe[kennung.Hinweis, *kennung.Hinweis]

func (c *CwdFiles) tryZettel(d string, a string, p string) (err error) {
	// kopf := filepath.Base(d)

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
	t.Kennung = h

	if path.Ext(a) == c.erworben.GetZettelFileExtension() {
		t.FDs.Objekte.Path = p
	} else {
		t.FDs.Akte.Path = p
	}

	return c.Zettelen.Add(t)
}

func (c CwdFiles) hinweisFromPath(p string) (h kennung.Hinweis, err error) {
	parts := strings.Split(p, string(filepath.Separator))

	switch len(parts) {
	case 0:
		fallthrough

	case 1:
		err = errors.Errorf("not enough parts: %q", parts)
		return

	default:
		parts = parts[len(parts)-2:]

	case 2:
		break
	}

	p = strings.Join(parts, string(filepath.Separator))

	p1 := p
	ext := path.Ext(p)

	if len(ext) != 0 {
		p1 = p[:len(p)-len(ext)]
	}

	if err = h.Set(p1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
