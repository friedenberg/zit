package verzeichnisse

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

type State int

const (
	StateUnread = State(iota)
	StateRead
	StateChanged
)

type zettelenPageWithState struct {
	path string
	ioFactory
	zettelenPage
	State
}

func makeZettelenPage(iof ioFactory, path string) (p *zettelenPageWithState) {
	p = &zettelenPageWithState{
		ioFactory: iof,
		path:      path,
		zettelenPage: zettelenPage{
			existing: make([]*Zettel, 0),
			added:    make([]*Zettel, 0),
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

	if _, err = zp.WriteTo(w); err != nil {
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

	if _, err = zp.zettelenPage.Copy(r, w); err != nil {
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

	if _, err = zp.ReadFrom(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	zp.State = StateRead

	return
}
