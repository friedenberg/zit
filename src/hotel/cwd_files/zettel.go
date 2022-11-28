package cwd_files

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
)

type CwdZettel struct {
	hinweis.Hinweis
	Zettel, Akte File
}

func (c *CwdFiles) tryZettel(d string, a string, p string) (err error) {
	var h hinweis.Hinweis

	kopf := filepath.Base(d)

	if h, err = c.hinweisFromPath(path.Join(kopf, a)); err != nil {
		err = errors.Wrap(err)
		return
	}

	var zcw CwdZettel
	ok := false

	if zcw, ok = c.Zettelen[h.String()]; !ok {
		zcw.Hinweis = h
	}

	//TODO read zettels
	if path.Ext(a) == c.konfig.Transacted.Objekte.GetZettelFileExtension() {
		zcw.Zettel.Path = p
	} else {
		zcw.Akte.Path = p
	}

	c.Zettelen[h.String()] = zcw

	return
}

func (c CwdFiles) hinweisFromPath(p string) (h hinweis.Hinweis, err error) {
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
