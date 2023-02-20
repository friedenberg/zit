package cwd

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type CwdZettel struct {
	kennung.Hinweis
	Zettel, Akte kennung.FD
}

func (c *CwdFiles) tryZettel(d string, a string, p string) (err error) {
	var h kennung.Hinweis

	kopf := filepath.Base(d)

	if h, err = c.hinweisFromPath(path.Join(kopf, a)); err != nil {
		err = errors.Wrap(err)
		return
	}

	var zcw *zettel.External
	ok := false

	if zcw, ok = c.Zettelen[h]; !ok {
		zcw = &zettel.External{
			Sku: sku.External[kennung.Hinweis, *kennung.Hinweis]{
				Kennung: h,
			},
		}
	}

	errors.TodoP3("read zettels")
	if path.Ext(a) == c.erworben.GetZettelFileExtension() {
		zcw.Sku.ObjekteFD.Path = p
	} else {
		zcw.Sku.AkteFD.Path = p
	}

	c.Zettelen[h] = zcw

	return
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
