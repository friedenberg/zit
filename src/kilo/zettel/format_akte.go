package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

type Akte struct {
	Toml        bool
	AkteFactory gattung.AkteIOFactory
}

func (f Akte) WriteTo(c FormatContextWrite) (n int64, err error) {
	var r io.ReadCloser

	sb := c.Zettel.Akte

	if r, err = f.AkteFactory.AkteReader(sb); err != nil {
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
