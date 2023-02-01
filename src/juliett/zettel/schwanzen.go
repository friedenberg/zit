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
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type Schwanzen struct {
	lock      *sync.RWMutex
	hinweisen map[kennung.Hinweis]sku.Transacted[kennung.Hinweis, *kennung.Hinweis]
}

func MakeSchwanzen() *Schwanzen {
	return &Schwanzen{
		lock:      &sync.RWMutex{},
		hinweisen: make(map[kennung.Hinweis]sku.Transacted[kennung.Hinweis, *kennung.Hinweis]),
	}
}

func (zws *Schwanzen) Less(zt *Transacted) (ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	var t sku.Transacted[kennung.Hinweis, *kennung.Hinweis]

	t, ok = zws.hinweisen[zt.Sku.Kennung]

	switch {
	case !ok:
		fallthrough

	case zt.Sku.Less(&t):
		ok = true
	}

	return
}

func (zws *Schwanzen) Get(h kennung.Hinweis) (t ts.Time, ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	var o sku.Transacted[kennung.Hinweis, *kennung.Hinweis]

	o, ok = zws.hinweisen[h]

	t = o.Schwanz

	return
}

func (zws *Schwanzen) Set(z *Transacted) (ok bool) {
	zws.lock.Lock()
	defer zws.lock.Unlock()

	h := z.Sku.Kennung
	o := z.Sku
	t1, _ := zws.hinweisen[h]

	if t1.Less(&o) {
		errors.Log().Printf("updating schwanzen %s -> %s", t1, o)
		zws.hinweisen[h] = o
		ok = true
	} else if t1.Equals(o) {
		ok = true
	} else {
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

	m := make(map[kennung.Hinweis]Sku)

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
