package zettel

import (
	"sync"

	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/kennung_index"
	"github.com/friedenberg/zit/src/india/transacted"
)

// TODO-P3 move to collections
type Schwanzen struct {
	lock         *sync.RWMutex
	hinweisen    map[kennung.Hinweis]transacted.Zettel
	etikettIndex kennung_index.EtikettIndex
}

func MakeSchwanzen(ei kennung_index.EtikettIndex) *Schwanzen {
	return &Schwanzen{
		lock:         &sync.RWMutex{},
		hinweisen:    make(map[kennung.Hinweis]transacted.Zettel),
		etikettIndex: ei,
	}
}

func (zws *Schwanzen) Less(zt *transacted.Zettel) (ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	t, ok := zws.hinweisen[zt.GetKennung()]

	switch {
	case !ok:
		fallthrough

	case zt.Less(t):
		ok = true
	}

	return
}

func (zws *Schwanzen) Get(h kennung.Hinweis) (t kennung.Tai, ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	o, ok := zws.hinweisen[h]

	if ok {
		t = o.GetTai()
	}

	return
}

func (zws *Schwanzen) Set(z *transacted.Zettel, flush bool) (ok bool) {
	zws.lock.Lock()
	defer zws.lock.Unlock()

	h := z.GetKennung()
	t1, found := zws.hinweisen[h]

	switch {
	case !found:
		fallthrough

	case t1.Less(*z):
		zws.hinweisen[h] = *z
		ok = true

	case t1.Metadatei.EqualsSansTai(z.Metadatei):
		zws.etikettIndex.Add(z.GetMetadatei().Etiketten)

		ok = flush && t1.GetTai().Equals(z.GetTai())

	default:
		zws.etikettIndex.AddEtikettSet(
			t1.GetMetadatei().Etiketten,
			z.GetMetadatei().Etiketten,
		)
	}

	return
}

func (zws *Schwanzen) ShouldAddVerzeichnisse(
	z *transacted.Zettel) (err error) {
	if ok := zws.Set(z, false); !ok {
		err = collections.MakeErrStopIteration()
		return
	}

	return
}

func (zws *Schwanzen) ShouldFlushVerzeichnisse(
	z *transacted.Zettel) (err error) {
	if ok := zws.Set(z, true); !ok {
		err = collections.MakeErrStopIteration()
		return
	}

	return
}