package zettel_verzeichnisse

import (
	"github.com/friedenberg/zit/src/delta/collections"
)

type Pool = collections.Pool[Zettel]

// type Pool struct {
// 	*collections.Pool[Zettel]
// }

// func MakePool() Pool {
// 	return Pool{
// 		Pool: collections.MakePool[Zettel](),
// 	}
// }

// func (ip Pool) WriteZettelVerzeichnisse(z *Zettel) (err error) {
// 	if z == nil {
// 		err = io.EOF
// 		return
// 	}

// 	ip.Put(z)
// 	return
// }

// func (p Pool) MakeZettel(
// 	tz zettel.Transacted,
// ) (z *Zettel) {
// 	z = p.Get()
// 	z.Transacted = tz
// 	z.EtikettenExpandedSorted = kennung.Expanded(tz.Objekte.Etiketten).SortedString()
// 	z.EtikettenSorted = tz.Objekte.Etiketten.SortedString()

// 	return
// }
