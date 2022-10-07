package store_verzeichnisse

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

type SchwanzEntryFunc func(z zettel_transacted.Zettel) (key string, value string)

type zettelenPageWithState struct {
	path string
	ioFactory
	zettelenPage
	zettelenPageIndex
	schwanzEntryFunc SchwanzEntryFunc
	State
}

func makeZettelenPage(
	iof ioFactory,
	path string,
	f SchwanzEntryFunc,
) (p *zettelenPageWithState) {
	p = &zettelenPageWithState{
		ioFactory: iof,
		path:      path,
		zettelenPage: zettelenPage{
			existing: make([]*Zettel, 0),
			added:    make([]*Zettel, 0),
		},
		schwanzEntryFunc: f,
		zettelenPageIndex: zettelenPageIndex{
			self: make(map[string]string),
		},
	}

	return
}

func (zp *zettelenPageWithState) Add(z *Zettel) (err error) {
	if err = zp.ReadAll(); err != nil {
		err = errors.Wrap(err)
		return
	}

	zp.State = StateChanged
	zp.added = append(zp.added, z)

	if z.PageSelection.Reason == PageSelectionReasonHinweis {
		key, value := zp.schwanzEntryFunc(z.Transacted)

		zp.self[key] = value
	}

	return
}

func (zpi *zettelenPageWithState) IsSchwanz(
	z zettel_transacted.Zettel,
) (ok bool, err error) {
	if err = zpi.ReadAll(); err != nil {
		err = errors.Wrap(err)
		return
	}

	key, value := zpi.schwanzEntryFunc(z)

	var value1 string

	value1, ok = zpi.self[key]

	switch {
	case !ok:
		return

	case value1 != value:
		ok = false

	default:
		ok = true
	}

	return
}

func (zp *zettelenPageWithState) ReadHinweis(
	h hinweis.Hinweis,
) (tz zettel_transacted.Zettel, err error) {
	if err = zp.ReadAll(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var z *Zettel

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
	if zp.State <= StateRead {
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

	if _, err = zp.zettelenPageIndex.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = zp.zettelenPage.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (zp *zettelenPageWithState) WriteZettelenTo(
	w writer,
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

	if _, err = zp.zettelenPageIndex.ReadFrom(r1); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = zp.zettelenPage.Copy(r1, w); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (zp *zettelenPageWithState) ReadAll() (err error) {
	if zp.State >= StateRead {
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

	if _, err = zp.zettelenPageIndex.ReadFrom(r1); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = zp.zettelenPage.ReadFrom(r1); err != nil {
		err = errors.Wrap(err)
		return
	}

	zp.State = StateRead

	return
}
