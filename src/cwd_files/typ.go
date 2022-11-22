package cwd_files

import (
	"path"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/typ"
)

func (c *CwdFiles) tryTyp(d string, a string, p string) (err error) {
	var h typ.Kennung

	ext := filepath.Ext(a)
	pathMinusExt := path.Base(a)[:len(ext)]

	if err = h.Set(pathMinusExt); err != nil {
		err = errors.Wrap(err)
		return
	}

	var t CwdTyp

	ok := false

	if t, ok = c.Typen[t.String()]; !ok {
		t.Kennung = h
	}

	t.Path = p
	c.Typen[h.String()] = t

	return
}
