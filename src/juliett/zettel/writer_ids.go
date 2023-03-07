package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
)

func init() {
	errors.TodoP2("move away from this and replace with compiled filter")
}

type WriterIds struct {
	Filter kennung.Filter
}

func (w WriterIds) WriteTransactedLike(maybeZ objekte.TransactedLike) (err error) {
	if z, ok := maybeZ.(*Transacted); ok {
		return w.WriteZettelTransacted(z)
	}

	return
}

func (w WriterIds) WriteZettelTransacted(z *Transacted) (err error) {
	return w.Filter.Include(z)
}
