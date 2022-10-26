package zettel_verzeichnisse

import (
	"sync"

	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
)

type Pool struct {
	inner *sync.Pool
}

func MakePool() *Pool {
	return &Pool{
		inner: &sync.Pool{
			New: func() interface{} {
				return &Zettel{}
			},
		},
	}
}

func (ip Pool) Get() *Zettel {
	return ip.inner.Get().(*Zettel)
}

func (ip Pool) Put(i *Zettel) {
	if i == nil {
		return
	}

	i.Reset(nil)
	ip.inner.Put(i)
}

func (ip Pool) WriteZettelVerzeichnisse(z *Zettel) (err error) {
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
