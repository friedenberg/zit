package store_verzeichnisse

import (
	"bufio"
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/zettel_verzeichnisse"
)

type zettelenPageWithState struct {
	sync.Locker
	pageId
	ioFactory
	zettelenPage
	State
}

func makeZettelenPage(
	iof ioFactory,
	pid pageId,
	pool *zettel_verzeichnisse.Pool,
) (p *zettelenPageWithState) {
	var flushFilter zettel_verzeichnisse.Writer
	flushFilter = zettel_verzeichnisse.WriterIdentity{}

	if zvwg, ok := iof.(ZettelVerzeichnisseWriterGetter); ok {
		flushFilter = zvwg.ZettelVerzeichnisseWriter(pid.index)
	}

	p = &zettelenPageWithState{
		Locker:    &sync.Mutex{},
		ioFactory: iof,
		pageId:    pid,
		zettelenPage: zettelenPage{
			pool:        pool,
			existing:    make([]*zettel_verzeichnisse.Zettel, 0),
			added:       make([]*zettel_verzeichnisse.Zettel, 0),
			flushFilter: flushFilter,
		},
	}

	return
}

func (zp *zettelenPageWithState) getState() State {
	zp.Lock()
	defer zp.Unlock()
	return zp.State
}

func (zp *zettelenPageWithState) setState(v State) {
	zp.Lock()
	defer zp.Unlock()
	zp.State = v
}

func (zp *zettelenPageWithState) Add(z *zettel_verzeichnisse.Zettel) (err error) {
	if err = zp.ReadAll(); err != nil {
		err = errors.Wrap(err)
		return
	}

	zp.setState(StateChanged)
	zp.added = append(zp.added, z)

	return
}

func (zp *zettelenPageWithState) ReadHinweis(
	h hinweis.Hinweis,
) (tz zettel_transacted.Zettel, err error) {
	if err = zp.ReadAll(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var z *zettel_verzeichnisse.Zettel

	for _, z1 := range zp.zettelenPage.existing {
		if z1.Transacted.Named.Hinweis.Equals(h) {
			z = z1
		}
	}

	for _, z1 := range zp.zettelenPage.added {
		if z1.Transacted.Named.Hinweis.Equals(h) {
			z = z1
		}
	}

	if z == nil {
		err = errors.Normalf("not found: %s", h)
	} else {
		tz = z.Transacted
	}

	return
}

func (zp *zettelenPageWithState) Flush() (err error) {
	state := zp.getState()
	if state <= StateRead {
		errors.Printf("no changes: %s", zp.path)
		return
	} else {
		errors.Printf("flushing: %s", zp.path)
	}

	var w io.WriteCloser

	if w, err = zp.WriteCloserVerzeichnisse(zp.path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.PanicIfError(w.Close)

	w1 := bufio.NewWriter(w)

	defer errors.PanicIfError(w1.Flush)

	if mpr, ok := zp.ioFactory.(PageHeader); ok {
		wt := mpr.PageHeaderWriterTo(zp.index)

		if _, err = wt.WriteTo(w1); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if _, err = zp.zettelenPage.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (zp *zettelenPageWithState) WriteZettelenTo(
	w zettel_verzeichnisse.Writer,
) (err error) {
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

	if mpr, ok := zp.ioFactory.(PageHeader); ok {
		rf := mpr.PageHeaderReaderFrom(zp.index)

		if _, err = rf.ReadFrom(r1); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if _, err = zp.zettelenPage.Copy(r1, w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (zp *zettelenPageWithState) ReadJustHeader() (err error) {
	var rf io.ReaderFrom

	if mpr, ok := zp.ioFactory.(PageHeader); ok {
		rf = mpr.PageHeaderReaderFrom(zp.index)
	} else {
		return
	}

	state := zp.getState()
	if state <= StateReadJustHeader {
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

	zp.setState(StateReadJustHeader)

	return
}

func (zp *zettelenPageWithState) ReadAll() (err error) {
	state := zp.getState()
	if state <= StateRead {
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

	if mpr, ok := zp.ioFactory.(PageHeader); ok {
		rf := mpr.PageHeaderReaderFrom(zp.index)

		if _, err = rf.ReadFrom(r1); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if _, err = zp.zettelenPage.ReadFrom(r1); err != nil {
		err = errors.Wrap(err)
		return
	}

	zp.setState(StateRead)

	return
}
