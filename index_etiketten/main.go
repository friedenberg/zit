package index_etiketten

import (
	"bufio"
	"encoding/gob"
	"io"

	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/delta/age_io"
)

type Index struct {
	path string
	age_io.ReadCloserFactory
	age_io.WriteCloserFactory
	etiketten  map[etikett.Etikett]int64
	didRead    bool
	hasChanges bool
}

type row struct {
	etikett.Etikett
	count int64
}

func newIndexEtiketten(
	p string,
	r age_io.ReadCloserFactory,
	w age_io.WriteCloserFactory,
) (i *Index, err error) {
	i = &Index{
		path:               p,
		ReadCloserFactory:  r,
		WriteCloserFactory: w,
		etiketten:          make(map[etikett.Etikett]int64),
	}

	return
}

func (i *Index) Flush() (err error) {
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

	for e, c := range i.etiketten {
		row := row{
			Etikett: e,
			count:   c,
		}

		logz.Print(row)

		if err = enc.Encode(row); err != nil {
			err = errors.Wrapped(err, "failed to write row: %s", row)
			return
		}
	}

	return
}

func (i *Index) readIfNecessary() (err error) {
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
		var r row

		if err = dec.Decode(&r); err != nil {
			if errors.IsEOF(err) {
				err = nil
				break
			} else {
				err = errors.Error(err)
				return
			}
		}

		i.etiketten[r.Etikett] = r.count
	}

	return
}

func (i *Index) Add(s etikett.Set) (err error) {
	if s.Len() == 0 {
		logz.Print("no etiketten to add")
		return
	}

	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	i.hasChanges = true

	for _, e := range s {
		logz.Printf("adding etiketten: %s", e)
		var c int64

		c, _ = i.etiketten[e]
		c += 1
		i.etiketten[e] = c
	}

	return
}

func (i *Index) Del(s etikett.Set) (err error) {
	if s.Len() == 0 {
		logz.Print("no etiketten to delete")
		return
	}

	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	i.hasChanges = true

	for _, e := range s {
		logz.Printf("removing etikett: %s", e)
		var c int64
		ok := false

		if c, ok = i.etiketten[e]; !ok {
			logz.Print(errors.Errorf("attempting to delete etikett that is already at 0"))
			return
		}

		c -= 1

		if c < 0 {
			delete(i.etiketten, e)
		} else {
			i.etiketten[e] = c
		}
	}

	return
}

func (i *Index) All() (es []etikett.Etikett, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	es = make([]etikett.Etikett, len(i.etiketten))

	n := 0

	for e, _ := range i.etiketten {
		es[n] = e
		n++
	}

	return
}
