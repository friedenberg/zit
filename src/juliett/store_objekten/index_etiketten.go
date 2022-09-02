package store_objekten

import (
	"bufio"
	"encoding/gob"
	"io"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/etikett"
)

type indexEtiketten struct {
	path string
	ioFactory
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
	f ioFactory,
) (i *indexEtiketten, err error) {
	i = &indexEtiketten{
		path:      p,
		ioFactory: f,
		etiketten: make(map[etikett.Etikett]int64),
	}

	return
}

func (i *indexEtiketten) Flush() (err error) {
	if !i.hasChanges {
		errors.Print("no changes")
		return
	}

	var w1 io.WriteCloser

	if w1, err = i.WriteCloserVerzeichnisse(i.path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.PanicIfError(w1.Close)

	w := bufio.NewWriter(w1)

	defer errors.PanicIfError(w.Flush)

	enc := gob.NewEncoder(w)

	for e, c := range i.etiketten {
		row := row{
			Etikett: e,
			count:   c,
		}

		if err = enc.Encode(row); err != nil {
			err = errors.Wrapf(err, "failed to write row: %s", row)
			return
		}
	}

	return
}

func (i *indexEtiketten) readIfNecessary() (err error) {
	if i.didRead {
		return
	}

	i.didRead = true

	var r1 io.ReadCloser

	if r1, err = i.ReadCloserVerzeichnisse(i.path); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
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
				err = errors.Wrap(err)
				return
			}
		}

		i.etiketten[r.Etikett] = r.count
	}

	return
}

func (i *indexEtiketten) add(s etikett.Set) (err error) {
	if s.Len() == 0 {
		errors.Print("no etiketten to add")
		return
	}

	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.hasChanges = true

	for _, e := range s {
		errors.Printf("adding etiketten: %s", e)
		var c int64

		c, _ = i.etiketten[e]
		c += 1
		i.etiketten[e] = c
	}

	return
}

func (i *indexEtiketten) del(s etikett.Set) (err error) {
	if s.Len() == 0 {
		errors.Print("no etiketten to delete")
		return
	}

	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.hasChanges = true

	for _, e := range s {
		errors.Printf("removing etikett: %s", e)
		var c int64
		ok := false

		if c, ok = i.etiketten[e]; !ok {
			errors.Print(errors.Errorf("attempting to delete etikett that is already at 0"))
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

func (i *indexEtiketten) allEtiketten() (es []etikett.Etikett, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	es = make([]etikett.Etikett, len(i.etiketten))

	n := 0

	for e, _ := range i.etiketten {
		es[n] = e
		n++
	}

	sort.Slice(
		es,
		func(i, j int) bool {
			return es[i].String() < es[j].String()
		},
	)

	return
}
