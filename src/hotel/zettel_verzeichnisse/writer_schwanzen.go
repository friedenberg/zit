package zettel_verzeichnisse

import (
	"bufio"
	"encoding/gob"
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

type schwanzValue struct {
	ts.Time
	sha.Sha
}

type WriterSchwanzen struct {
	lock      *sync.RWMutex
	hinweisen map[hinweis.Hinweis]schwanzValue
}

func MakeWriterSchwanzen() *WriterSchwanzen {
	return &WriterSchwanzen{
		lock:      &sync.RWMutex{},
		hinweisen: make(map[hinweis.Hinweis]schwanzValue),
	}
}

func (zws *WriterSchwanzen) Less(zt *zettel_transacted.Zettel) (ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	var t schwanzValue

	t, ok = zws.hinweisen[zt.Named.Hinweis]

	switch {
	case !ok:
		fallthrough

	case zt.Schwanz.Less(t.Time):
		ok = true
	}

	return
}

func (zws *WriterSchwanzen) Get(h hinweis.Hinweis) (t ts.Time, ok bool) {
	zws.lock.RLock()
	defer zws.lock.RUnlock()

	var sw schwanzValue

	sw, ok = zws.hinweisen[h]

	t = sw.Time

	return
}

func (zws *WriterSchwanzen) Set(z *zettel_transacted.Zettel) (ok bool) {
	zws.lock.Lock()
	defer zws.lock.Unlock()

	t := z.Schwanz
	h := z.Named.Hinweis
	sh := z.Named.Stored.Sha
	t1, _ := zws.hinweisen[h]

	if t1.Time.Equals(t) {
		if t1.Sha.Equals(sh) {
			ok = true
		} else {
			//This fixes an issue where some transactions have zettels appear more than
			//once
			ok = false
		}
	} else if t1.Less(t) {
		zws.hinweisen[h] = schwanzValue{Time: t, Sha: sh}
		ok = true
	} else {
		ok = false
	}

	return
}

func (zws *WriterSchwanzen) WriteZettelTransacted(
	z *zettel_transacted.Zettel,
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
	return zws.WriteZettelTransacted(&z.Transacted)
}

func (zws *WriterSchwanzen) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	m := make(map[hinweis.Hinweis]schwanzValue)

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

	defer errors.PanicIfError(w.Flush)

	enc := gob.NewEncoder(w)

	if err = enc.Encode(zws.hinweisen); err != nil {
		err = errors.Wrapf(err, "failed to write page index")
		return
	}

	return
}
