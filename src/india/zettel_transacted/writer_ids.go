package zettel_transacted

import (
	"io"

	"github.com/friedenberg/zit/src/hotel/zettel_named"
)

type WriterIds zettel_named.FilterIdSet

func (w WriterIds) NamedFilter() zettel_named.NamedFilter {
	return zettel_named.FilterIdSet(w)
}

func (w WriterIds) WriteZettelTransacted(z *Zettel) (err error) {
	if !w.NamedFilter().IncludeNamedZettel(&z.Named) {
		err = io.EOF
	}

	return
}
