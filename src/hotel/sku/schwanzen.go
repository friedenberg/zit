package sku

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/kennung_index"
)

type Schwanzen struct {
	lock         *sync.RWMutex
	hinweisen    map[string]Transacted
	etikettIndex kennung_index.EtikettIndex
	funcFlush    schnittstellen.FuncIter[*Transacted]
}

func MakeSchwanzen(
	ei kennung_index.EtikettIndex,
	funcFlush schnittstellen.FuncIter[*Transacted],
) *Schwanzen {
	return &Schwanzen{
		lock:         &sync.RWMutex{},
		hinweisen:    make(map[string]Transacted),
		etikettIndex: ei,
		funcFlush:    funcFlush,
	}
}

func (zws *Schwanzen) Less(zt *Transacted) (ok bool) {
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

func (zws *Schwanzen) Set(z *Transacted, flush bool) (ok bool) {
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

	if err = zws.funcFlush(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
