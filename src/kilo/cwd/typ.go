package cwd

import (
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type Typ = sku.ExternalMaybe[kennung.Typ, *kennung.Typ]

func (c *CwdFiles) tryTyp(fi os.FileInfo) (err error) {
	var h kennung.Typ

	fd := kennung.FileInfo(fi)
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
