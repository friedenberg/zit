package store_verzeichnisse

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

type Page struct {
	sync.Locker
	pageId
	ioFactory
	pool        *collections.Pool[zettel.Transacted]
	added       zettel.HeapTransacted
	flushFilter collections.WriterFunc[*zettel.Transacted]
	State
}

func makeZettelenPage(
	iof ioFactory,
	pid pageId,
	pool *collections.Pool[zettel.Transacted],
	fff ZettelTransactedWriterGetter,
) (p *Page) {
	flushFilter := collections.MakeWriterNoop[*zettel.Transacted]()

	if fff != nil {
		flushFilter = fff.ZettelTransactedWriter(pid.index)
	}

	p = &Page{
		Locker:      &sync.Mutex{},
		ioFactory:   iof,
		pageId:      pid,
		pool:        pool,
		added:       zettel.MakeHeapTransacted(),
		flushFilter: flushFilter,
	}

	return
}

func (zp *Page) getState() State {
	zp.Lock()
	defer zp.Unlock()
	return zp.State
}

func (zp *Page) setState(v State) {
	zp.Lock()
	defer zp.Unlock()
	zp.State = v
}

func (zp *Page) Add(z *zettel.Transacted) (err error) {
	if z == nil {
		err = errors.Errorf("trying to add nil zettel.Verzeichnisse")
		return
	}

	if err = zp.flushFilter(z); err != nil {
		if errors.Is(err, collections.ErrStopIteration) {
			errors.Log().Printf("eliding %s", z.Kennung())
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	if z == nil {
		err = errors.Errorf("trying to add nil zettel.Zettel")
		return
	}

	zp.added.Add(*z)
	zp.setState(StateChanged)

	return
}

func (zp *Page) Flush() (err error) {
	if zp.getState() < StateChanged {
		return
	}

	errors.Log().Printf("flushing page: %s", zp.path)

	var w io.WriteCloser

	if w, err = zp.WriteCloserVerzeichnisse(zp.path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, w.Close)

	w1 := bufio.NewWriter(w)

	defer errors.Deferred(&err, w1.Flush)

	if mpr, ok := zp.ioFactory.(PageHeader); ok {
		wt := mpr.PageHeaderWriterTo(zp.index)

		if _, err = wt.WriteTo(w1); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if _, err = zp.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	zp.added.Reset()

	return
}

func (zp *Page) WriteZettelenTo(
	w collections.WriterFunc[*zettel.Transacted],
) (err error) {
	var r io.ReadCloser

	if r, err = zp.ReadCloserVerzeichnisse(zp.path); err != nil {
		if errors.IsNotExist(err) {
			r = ioutil.NopCloser(bytes.NewReader(nil))
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	defer errors.Deferred(&err, r.Close)

	r1 := bufio.NewReader(r)

	if mpr, ok := zp.ioFactory.(PageHeader); ok {
		rf := mpr.PageHeaderReaderFrom(zp.index)

		if _, err = rf.ReadFrom(r1); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if _, err = zp.Copy(r1, w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (zp *Page) ReadJustHeader() (err error) {
	var rf io.ReaderFrom

	if mpr, ok := zp.ioFactory.(PageHeader); ok {
		rf = mpr.PageHeaderReaderFrom(zp.index)
	} else {
		return
	}

	state := zp.getState()
	if state <= StateReadHeader {
		errors.Log().Printf("already read %s", zp.path)
		return
	} else {
		errors.Log().Printf("reading: %s", zp.path)
	}

	var r io.ReadCloser

	if r, err = zp.ReadCloserVerzeichnisse(zp.path); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.Deferred(&err, r.Close)

	r1 := bufio.NewReader(r)

	if _, err = rf.ReadFrom(r1); err != nil {
		err = errors.Wrap(err)
		return
	}

	zp.setState(StateReadHeader)

	return
}

func (zp *Page) Copy(
	r1 io.Reader,
	w collections.WriterFunc[*zettel.Transacted],
) (n int64, err error) {
	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	defer func() {
		zp.added.Restore()
	}()

	for {
		var tz *zettel.Transacted

		tz = zp.pool.Get()

		if err = dec.Decode(tz); err != nil {
			if errors.IsEOF(err) {
				err = nil
				break
			} else {
				err = errors.Wrap(err)
				return
			}
		}

	LOOP:
		for {
			peeked, ok := zp.added.PeekPtr()

			switch {
			case !ok:
				break LOOP

			case peeked.Equals(tz):
				zp.added.PopAndSave()
				break

			case !peeked.Less(*tz):
				break LOOP
			}

			popped, _ := zp.added.PopAndSave()

			if err = w(&popped); err != nil {
				if collections.IsStopIteration(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}
		}

		if err = w(tz); err != nil {
			if collections.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	for {
		popped, ok := zp.added.PopAndSave()

		if !ok {
			break
		}

		if err = w(&popped); err != nil {
			if errors.Is(err, collections.ErrStopIteration) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (zp *Page) WriteTo(w1 io.Writer) (n int64, err error) {
	w := bufio.NewWriter(w1)

	defer errors.Deferred(&err, w.Flush)

	if err = zp.WriteZettelenTo(
		collections.MakeChain(
			zp.flushFilter,
			zettel.MakeWriterGobEncoder(w).WriteZettelVerzeichnisse,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
