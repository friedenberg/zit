package cwd

import (
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type Kasten = sku.ExternalMaybe[kennung.Kasten, *kennung.Kasten]

func (c *CwdFiles) tryKasten(fi os.FileInfo, dir string) (err error) {
	var h kennung.Kasten
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

	t, _ := c.Kisten.Get(h.String())

	t.Kennung = h
	t.FDs.Objekte = fd
	return c.Kisten.Add(t)
}
