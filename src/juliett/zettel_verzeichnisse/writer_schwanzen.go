package zettel_verzeichnisse

import (
	"bufio"
	"encoding/gob"
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/sku"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
)

type WriterSchwanzen struct {
	lock      *sync.RWMutex
	hinweisen map[hinweis.Hinweis]sku.Transacted
}

func MakeWriterSchwanzen() *WriterSchwanzen {
	return &WriterSchwanzen{
		lock:      &sync.RWMutex{},
		hinweisen: make(map[hinweis.Hinweis]sku.Transacted),
	}
}

func (zws *WriterSchwanzen) Less(zt *zettel_transacted.Transacted) (ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	var t sku.Transacted

	t, ok = zws.hinweisen[zt.Named.Kennung]

	switch {
	case !ok:
		fallthrough

	case zt.SkuTransacted().Less(t):
		ok = true
	}

	return
}

func (zws *WriterSchwanzen) Get(h hinweis.Hinweis) (t ts.Time, ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	var o sku.Transacted

	o, ok = zws.hinweisen[h]

	t = o.Schwanz

	return
}

func (zws *WriterSchwanzen) Set(z *zettel_transacted.Transacted) (ok bool) {
	zws.lock.Lock()
	defer zws.lock.Unlock()

	h := z.Named.Kennung
	o := z.SkuTransacted()
	t1, _ := zws.hinweisen[h]

	if t1.Less(o) {
		zws.hinweisen[h] = o
		ok = true
	} else if t1.Equals(o) {
		ok = true
	} else {
		ok = false
	}

	return
}

func (zws *WriterSchwanzen) WriteZettelTransacted(
	z *zettel_transacted.Transacted,
) (err error) {
	if ok := zws.Set(z); !ok {
		err = io.EOF
		return
	}

	return
}

func (zws *WriterSchwanzen) WriteZettelVerzeichnisse(
	z *Zettel,
) (err error) {
	err = zws.WriteZettelTransacted(&z.Transacted)

	return
}

func (zws *WriterSchwanzen) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	m := make(map[hinweis.Hinweis]sku.Transacted)

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
