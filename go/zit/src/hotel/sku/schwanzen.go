package sku

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_id_index"
)

type Schwanzen struct {
	lock               sync.RWMutex
	object_id_provider map[string]*Transacted
	etikettIndex       object_id_index.EtikettIndexMutation
	funcFlush          interfaces.FuncIter[*Transacted]
}

func MakeSchwanzen(
	ei object_id_index.EtikettIndexMutation,
	funcFlush interfaces.FuncIter[*Transacted],
) (s *Schwanzen) {
	s = &Schwanzen{}
	s.Initialize(ei, funcFlush)
	return
}

func (s *Schwanzen) Initialize(
	ei object_id_index.EtikettIndexMutation,
	funcFlush interfaces.FuncIter[*Transacted],
) {
	s.object_id_provider = make(map[string]*Transacted)
	s.etikettIndex = ei
	s.funcFlush = funcFlush
}

func (zws *Schwanzen) Less(zt *Transacted) (ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	t, ok := zws.object_id_provider[zt.GetKennung().String()]

	switch {
	case !ok:
		fallthrough

	case TransactedLessor.LessPtr(zt, t):
		ok = true
	}

	return
}

func (zws *Schwanzen) Get(h ids.IdLike) (t ids.Tai, ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	o, ok := zws.object_id_provider[h.String()]

	if ok {
		t = o.GetTai()
	}

	return
}

func (zws *Schwanzen) Set(z *Transacted, flush bool) (ok bool) {
	zws.lock.Lock()
	defer zws.lock.Unlock()

	h := z.GetKennung()
	t1, found := zws.object_id_provider[h.String()]

	switch {
	case !found:
		fallthrough

	case TransactedLessor.LessPtr(t1, z):
		zws.object_id_provider[h.String()] = z
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
