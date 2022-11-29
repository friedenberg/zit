package zettel_verzeichnisse

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/india/zettel"
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
	tz zettel.Transacted,
) (z *Zettel) {
	z = p.Get()
	z.Transacted = tz
	z.EtikettenExpandedSorted = kennung.Expanded(tz.Named.Stored.Objekte.Etiketten).SortedString()
	z.EtikettenSorted = tz.Named.Stored.Objekte.Etiketten.SortedString()

	return
}
