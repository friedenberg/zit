package zettel_verzeichnisse

import (
	"sync"

	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

type PoolGetter interface {
	ZettelVerzeichnissePool() *Pool
}

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
	i.Reset()
	ip.inner.Put(i)
}

func (p Pool) MakeZettel(
	tz zettel_transacted.Zettel,
) (z *Zettel) {
	z = p.Get()
	z.Transacted = tz
	z.EtikettenExpandedSorted = tz.Named.Stored.Zettel.Etiketten.Expanded().SortedString()
	z.EtikettenSorted = tz.Named.Stored.Zettel.Etiketten.SortedString()

	return
}
