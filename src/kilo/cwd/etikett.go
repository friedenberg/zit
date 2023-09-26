package cwd

import (
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type Etikett = sku.ExternalMaybe

func (c *CwdFiles) tryEtikett(fi os.FileInfo, dir string) (err error) {
	var h kennung.Etikett
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

	t, _ := c.Etiketten.Get(h.String())

	t.Kennung = kennung.Kennung2{KennungPtr: &h}
	t.FDs.Objekte = fd
	return c.Etiketten.Add(t)
}
