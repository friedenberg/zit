package zettel

import (
	"sync"

	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
)

// TODO-P3 move to collections
type Schwanzen struct {
	lock         *sync.RWMutex
	hinweisen    map[kennung.Hinweis]Transacted
	etikettIndex kennung_index.EtikettIndex
}

func MakeSchwanzen(ei kennung_index.EtikettIndex) *Schwanzen {
	return &Schwanzen{
		lock:         &sync.RWMutex{},
		hinweisen:    make(map[kennung.Hinweis]Transacted),
		etikettIndex: ei,
	}
}

func (zws *Schwanzen) Less(zt *Transacted) (ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	t, ok := zws.hinweisen[zt.Sku.Kennung]

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
		t = o.Sku.GetTai()
	}

	return
}

func (zws *Schwanzen) Set(z *Transacted, flush bool) (ok bool) {
	zws.lock.Lock()
	defer zws.lock.Unlock()

	h := z.Sku.Kennung
	t1, found := zws.hinweisen[h]

	switch {
	case !found:
		fallthrough

	case t1.Less(*z):
		zws.hinweisen[h] = *z
		ok = true

	case t1.Sku.Equals(z.Sku):
		zws.etikettIndex.Add(z.GetMetadatei().Etiketten)

		ok = flush

	default:
		zws.etikettIndex.AddEtikettSet(
			t1.GetMetadatei().Etiketten,
			z.GetMetadatei().Etiketten,
		)
	}

	return
}

func (zws *Schwanzen) ShouldAddVerzeichnisse(
	z *Transacted,
) (err error) {
	if ok := zws.Set(z, false); !ok {
		err = collections.MakeErrStopIteration()
		return
	}

	return
}

func (zws *Schwanzen) ShouldFlushVerzeichnisse(
	z *Transacted,
) (err error) {
	if ok := zws.Set(z, true); !ok {
		err = collections.MakeErrStopIteration()
		return
	}

	return
}
