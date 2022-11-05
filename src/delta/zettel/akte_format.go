package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Akte struct {
	Toml bool
}

func (f Akte) WriteTo(c FormatContextWrite) (n int64, err error) {
	var r io.ReadCloser

	sb := c.Zettel.Akte

	if r, err = c.AkteReader(sb); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, r.Close)

	if _, err = io.Copy(c.Out, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *Akte) ReadFrom(c *FormatContextRead) (n int64, err error) {
	return
}
