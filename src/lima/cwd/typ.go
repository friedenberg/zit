package cwd

import (
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type Typ = sku.ExternalMaybe[kennung.Typ, *kennung.Typ]

func (c *CwdFiles) tryTyp(fi os.FileInfo, dir string) (err error) {
	var h kennung.Typ
	var fd kennung.FD

	if fd, err = kennung.FileInfo(fi, dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	pathMinusExt := strings.ToLower(fd.FileNameSansExt())

	if err = h.Set(pathMinusExt); err != nil {
		err = errors.Wrap(err)
		return
	}

	t, _ := c.Typen.Get(h.String())

	t.Kennung = h
	t.FDs.Objekte = fd
	return c.Typen.Add(t)
}
