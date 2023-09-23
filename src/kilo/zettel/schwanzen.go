package zettel

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/kennung_index"
	"github.com/friedenberg/zit/src/hotel/sku"
)

// TODO-P3 move to collections
type Schwanzen struct {
	lock         *sync.RWMutex
	hinweisen    map[string]sku.Transacted
	etikettIndex kennung_index.EtikettIndex
	funcFlush    schnittstellen.FuncIter[*sku.Transacted]
}

func MakeSchwanzen(
	ei kennung_index.EtikettIndex,
	funcFlush schnittstellen.FuncIter[*sku.Transacted],
) *Schwanzen {
	return &Schwanzen{
		lock:         &sync.RWMutex{},
		hinweisen:    make(map[string]sku.Transacted),
		etikettIndex: ei,
		funcFlush:    funcFlush,
	}
}

func (zws *Schwanzen) Less(zt *sku.Transacted) (ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	t, ok := zws.hinweisen[zt.GetKennung().String()]

	switch {
	case !ok:
		fallthrough

	case zt.Less(t):
		ok = true
	}

	return
}

func (zws *Schwanzen) Get(h kennung.Kennung) (t kennung.Tai, ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	o, ok := zws.hinweisen[h.String()]

	if ok {
		t = o.GetTai()
	}

	return
}

func (zws *Schwanzen) Set(z *sku.Transacted, flush bool) (ok bool) {
	zws.lock.Lock()
	defer zws.lock.Unlock()

	h := z.GetKennung()
	t1, found := zws.hinweisen[h.String()]

	switch {
	case !found:
		fallthrough

	case t1.Less(*z):
		zws.hinweisen[h.String()] = *z
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
	z *sku.Transacted,
) (err error) {
	if ok := zws.Set(z, false); !ok {
		err = collections.MakeErrStopIteration()
		return
	}

	return
}

func (zws *Schwanzen) ShouldFlushVerzeichnisse(
	z *sku.Transacted,
) (err error) {
	if ok := zws.Set(z, true); !ok {
		err = collections.MakeErrStopIteration()
		return
	}

	if err = zws.funcFlush(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
