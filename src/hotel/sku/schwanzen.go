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
	lock         sync.RWMutex
	hinweisen    map[string]*Transacted
	etikettIndex kennung_index.EtikettIndexMutation
	funcFlush    schnittstellen.FuncIter[*Transacted]
}

func MakeSchwanzen(
	ei kennung_index.EtikettIndexMutation,
	funcFlush schnittstellen.FuncIter[*Transacted],
) (s *Schwanzen) {
	s = &Schwanzen{}
	s.Initialize(ei, funcFlush)
  return
}

func (s *Schwanzen) Initialize(
	ei kennung_index.EtikettIndexMutation,
	funcFlush schnittstellen.FuncIter[*Transacted],
) {
	s.hinweisen = make(map[string]*Transacted)
	s.etikettIndex = ei
	s.funcFlush = funcFlush
}

func (zws *Schwanzen) Less(zt *Transacted) (ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	t, ok := zws.hinweisen[zt.GetKennung().String()]

	switch {
	case !ok:
		fallthrough

	case TransactedLessor.LessPtr(zt, t):
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

	case TransactedLessor.LessPtr(t1, z):
		zws.hinweisen[h.String()] = z
		ok = true

	case t1.Metadatei.EqualsSansTai(&z.Metadatei):
		zws.etikettIndex.Add(z.Metadatei.GetEtiketten())
		ok = flush && t1.GetTai().Equals(z.GetTai())

	default:
		zws.etikettIndex.AddEtikettSet(
			t1.Metadatei.GetEtiketten(),
			z.Metadatei.GetEtiketten(),
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
