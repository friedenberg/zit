package objekten

import (
	"bufio"
	"encoding/gob"
	"io"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/objekte"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type indexZettelen struct {
	path string
	objekte.ReadCloserFactory
	objekte.WriteCloserFactory
	zettelen   map[hinweis.Hinweis]stored_zettel.Transacted
	didRead    bool
	hasChanges bool
}

func newIndexZettelen(
	p string,
	r objekte.ReadCloserFactory,
	w objekte.WriteCloserFactory,
) (i *indexZettelen, err error) {
	i = &indexZettelen{
		path:               p,
		ReadCloserFactory:  r,
		WriteCloserFactory: w,
		zettelen:           make(map[hinweis.Hinweis]stored_zettel.Transacted),
	}

	return
}

func (i *indexZettelen) Flush() (err error) {
	if !i.hasChanges {
		logz.Print("no changes")
		return
	}

	var w1 io.WriteCloser

	if w1, err = i.WriteCloser(i.path); err != nil {
		err = errors.Error(err)
		return
	}

	defer stdprinter.PanicIfError(w1.Close)

	w := bufio.NewWriter(w1)

	defer stdprinter.PanicIfError(w.Flush)

	enc := gob.NewEncoder(w)

	for h, tz := range i.zettelen {
		logz.Print(h)

		if err = enc.Encode(tz); err != nil {
			err = errors.Wrapped(err, "failed to write zettel: [%s %s]", h, tz.Sha)
			return
		}
	}

	return
}

func (i *indexZettelen) readIfNecessary() (err error) {
	if i.didRead {
		return
	}

	i.didRead = true

	var r1 io.ReadCloser

	if r1, err = i.ReadCloser(i.path); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Error(err)
		}
		return
	}

	defer r1.Close()

	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	for {
		var tz stored_zettel.Transacted

		if err = dec.Decode(&tz); err != nil {
			if errors.IsEOF(err) {
				err = nil
				break
			} else {
				err = errors.Error(err)
				return
			}
		}

		i.zettelen[tz.Hinweis] = tz
	}

	return
}

func (i *indexZettelen) Add(tz stored_zettel.Transacted) (err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	i.hasChanges = true

	i.zettelen[tz.Hinweis] = tz

	return
}

func (i *indexZettelen) Read(h hinweis.Hinweis) (tz stored_zettel.Transacted, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	ok := false

	if tz, ok = i.zettelen[h]; !ok {
		err = ErrNotFound{Id: h}
		return
	}

	return
}

func (i *indexZettelen) allTransacted() (tz map[hinweis.Hinweis]stored_zettel.Transacted, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	tz = make(map[hinweis.Hinweis]stored_zettel.Transacted)

	for h, z := range i.zettelen {
		tz[h] = z
	}

	return
}
