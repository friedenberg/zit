package zettel

import (
	"bufio"
	"encoding/gob"
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/sku"
)

type WriterSchwanzen struct {
	lock      *sync.RWMutex
	hinweisen map[hinweis.Hinweis]sku.Transacted[hinweis.Hinweis, *hinweis.Hinweis]
}

func MakeWriterSchwanzen() *WriterSchwanzen {
	return &WriterSchwanzen{
		lock:      &sync.RWMutex{},
		hinweisen: make(map[hinweis.Hinweis]sku.Transacted[hinweis.Hinweis, *hinweis.Hinweis]),
	}
}

func (zws *WriterSchwanzen) Less(zt *Transacted) (ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	var t sku.Transacted[hinweis.Hinweis, *hinweis.Hinweis]

	t, ok = zws.hinweisen[zt.Sku.Kennung]

	switch {
	case !ok:
		fallthrough

	case zt.Sku.Less(&t):
		ok = true
	}

	return
}

func (zws *WriterSchwanzen) Get(h hinweis.Hinweis) (t ts.Time, ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	var o sku.Transacted[hinweis.Hinweis, *hinweis.Hinweis]

	o, ok = zws.hinweisen[h]

	t = o.Schwanz

	return
}

func (zws *WriterSchwanzen) Set(z *Transacted) (ok bool) {
	zws.lock.Lock()
	defer zws.lock.Unlock()

	h := z.Sku.Kennung
	o := z.Sku
	t1, _ := zws.hinweisen[h]

	if t1.Less(&o) {
		zws.hinweisen[h] = o
		ok = true
	} else if t1.Equals(&o) {
		ok = true
	} else {
		ok = false
	}

	return
}

func (zws *WriterSchwanzen) WriteZettelTransacted(
	z *Transacted,
) (err error) {
	if ok := zws.Set(z); !ok {
		err = io.EOF
		return
	}

	return
}

func (zws *WriterSchwanzen) WriteZettelVerzeichnisse(
	z *Verzeichnisse) (err error) {
	err = zws.WriteZettelTransacted(&z.Transacted)

	return
}

func (zws *WriterSchwanzen) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	m := make(map[hinweis.Hinweis]Sku)

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

func (zws WriterSchwanzen) WriteTo(w1 io.Writer) (n int64, err error) {
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