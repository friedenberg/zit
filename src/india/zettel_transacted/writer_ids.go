package zettel_transacted

import (
	"io"

	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

type WriterIds zettel.FilterIdSet

func (w WriterIds) NamedFilter() zettel.FilterIdSet {
	return zettel.FilterIdSet(w)
}

func (w WriterIds) WriteZettelTransacted(z *Zettel) (err error) {
	if !w.NamedFilter().IncludeNamedZettel(&z.Named) {
		err = io.EOF
	}

	return
}
