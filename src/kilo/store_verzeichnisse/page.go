package store_verzeichnisse

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/juliett/zettel_verzeichnisse"
)

type Page struct {
	sync.Locker
	pageId
	ioFactory
	pool        zettel_verzeichnisse.Pool
	added       []*zettel_verzeichnisse.Zettel
	flushFilter collections.WriterFunc[*zettel_verzeichnisse.Zettel]
	State
}

func makeZettelenPage(
	iof ioFactory,
	pid pageId,
	pool zettel_verzeichnisse.Pool,
) (p *Page) {
	flushFilter := collections.MakeWriterNoop[*zettel_verzeichnisse.Zettel]()

	if zvwg, ok := iof.(ZettelVerzeichnisseWriterGetter); ok {
		flushFilter = zvwg.ZettelVerzeichnisseWriter(pid.index)
	}

	p = &Page{
		Locker:      &sync.Mutex{},
		ioFactory:   iof,
		pageId:      pid,
		pool:        pool,
		added:       make([]*zettel_verzeichnisse.Zettel, 0),
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

func (zp *Page) Add(z *zettel_verzeichnisse.Zettel) (err error) {
	if err = zp.flushFilter(z); err != nil {
		if errors.IsEOF(err) {
			errors.Log().Printf("eliding %s", z.Transacted.Named.Hinweis)
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	zp.added = append(zp.added, z)
	zp.setState(StateChanged)

	return
}

func (zp *Page) Flush() (err error) {
	if zp.getState() < StateChanged {
		return
	}

	errors.Printf("flushing page: %s", zp.path)

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

	return
}

func (zp *Page) WriteZettelenTo(
	w collections.WriterFunc[*zettel_verzeichnisse.Zettel],
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

	defer r.Close()

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
		errors.Printf("already read %s", zp.path)
		return
	} else {
		errors.Printf("reading: %s", zp.path)
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

	defer r.Close()

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
	w collections.WriterFunc[*zettel_verzeichnisse.Zettel],
) (n int64, err error) {
	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	for {
		var tz *zettel_verzeichnisse.Zettel

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

		if err = w(tz); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	for _, z := range zp.added {
		z1 := zp.pool.Get()
		z1.Reset(z)

		if err = w(z1); err != nil {
			if errors.IsEOF(err) {
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
			zettel_verzeichnisse.MakeWriterGobEncoder(w).WriteZettelVerzeichnisse,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}