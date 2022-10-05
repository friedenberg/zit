package verzeichnisse

import (
	"bufio"
	"encoding/gob"
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
			innerSet: zettel_transacted.MakeSetHinweis(0),
		},
	}

	return
}

func (zp *zettelenPageWithState) Add(zt zettel_transacted.Zettel) (err error) {
	if err = zp.ReadAll(); err != nil {
		err = errors.Wrap(err)
		return
	}

	zp.State = StateChanged
	zp.innerSet.Add(zt)

	return
}

func (zp *zettelenPageWithState) ReadHinweis(
	h hinweis.Hinweis,
) (tz zettel_transacted.Zettel, err error) {
	if err = zp.ReadAll(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ok := false

	if tz, ok = zp.innerSet.Get(h); !ok {
		err = errors.Normalf("not found")
		return
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

func (zp *zettelenPageWithState) WriteZettelenTo(w zettel_transacted.Writer) (err error) {
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

type zettelenPage struct {
	innerSet zettel_transacted.Set
}

func (zp zettelenPage) Copy(r1 io.Reader, w zettel_transacted.Writer) (n int64, err error) {
	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	for {
		var tz zettel_transacted.Zettel

		if err = dec.Decode(&tz); err != nil {
			if errors.IsEOF(err) {
				err = nil
				break
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		if err = w.WriteZettelTransacted(tz); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (zp *zettelenPage) ReadFrom(r1 io.Reader) (n int64, err error) {
	return zp.Copy(r1, zp.innerSet)
}

func (zp *zettelenPage) WriteTo(w1 io.Writer) (n int64, err error) {
	w := bufio.NewWriter(w1)

	defer errors.PanicIfError(w.Flush)

	enc := gob.NewEncoder(w)

	err = zp.innerSet.Each(
		func(tz zettel_transacted.Zettel) (err error) {
			if err = enc.Encode(tz); err != nil {
				err = errors.Wrapf(err, "failed to write zettel: %s", tz.Named)
				return
			}

			return
		},
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
