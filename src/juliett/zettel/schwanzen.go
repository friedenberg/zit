package zettel

import (
	"bufio"
	"encoding/gob"
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/foxtrot/kennung_index"
)

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

func (zws *Schwanzen) Get(h kennung.Hinweis) (t ts.Time, ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	o, ok := zws.hinweisen[h]

	if ok {
		errors.TodoP4("switch to GetTime()")
		t = o.Sku.Schwanz
	}

	return
}

func (zws *Schwanzen) Set(z *Transacted) (ok bool) {
	zws.lock.Lock()
	defer zws.lock.Unlock()

	h := z.Sku.Kennung
	t1, found := zws.hinweisen[h]

	switch {
	case !found:
		zws.hinweisen[h] = *z
		ok = true

	case t1.Less(*z):
		zws.hinweisen[h] = *z
		ok = true

		errors.TodoP4("determine if comparing zettels rather than sku is ok")
	case t1.Sku.Equals(z.Sku):
		zws.etikettIndex.Add(z.Objekte.Etiketten)
		ok = true

	default:
		zws.etikettIndex.AddEtikettSet(t1.Objekte.Etiketten, z.Objekte.Etiketten)
		ok = false
	}

	return
}

func (zws *Schwanzen) WriteZettelTransacted(
	z *Transacted,
) (err error) {
	if ok := zws.Set(z); !ok {
		err = collections.MakeErrStopIteration()
		return
	}

	return
}

func (zws *Schwanzen) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	m := make(map[kennung.Hinweis]Transacted)

	if err = dec.Decode(&m); err != nil {
		if errors.IsEOF(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	zws.lock.Lock()
	defer zws.lock.Unlock()

	zws.hinweisen = m

	return
}

func (zws Schwanzen) WriteTo(w1 io.Writer) (n int64, err error) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	w := bufio.NewWriter(w1)

	defer errors.Deferred(&err, w.Flush)

	enc := gob.NewEncoder(w)

	if err = enc.Encode(zws.hinweisen); err != nil {
		err = errors.Wrapf(err, "failed to write page index")
		return
	}

	return
}
