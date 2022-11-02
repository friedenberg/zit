package zettel_verzeichnisse

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
)

type Pool struct {
	*collections.Pool[Zettel]
}

func MakePool() Pool {
	return Pool{
		Pool: collections.MakePool[Zettel](),
	}
}

func (ip Pool) WriteZettelVerzeichnisse(z *Zettel) (err error) {
	if z == nil {
		err = io.EOF
		return
	}

	ip.Put(z)
	return
}

func (p Pool) MakeZettel(
	tz zettel_transacted.Zettel,
) (z *Zettel) {
	z = p.Get()
	z.Transacted = tz
	z.EtikettenExpandedSorted = etikett.Expanded(tz.Named.Stored.Zettel.Etiketten).SortedString()
	z.EtikettenSorted = tz.Named.Stored.Zettel.Etiketten.SortedString()

	return
}
