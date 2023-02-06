package kennung_index

import (
	"bufio"
	"encoding/gob"
	"io"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/hinweis_index"
)

type EtikettIndex interface {
	GetAllEtiketten() ([]kennung.Etikett, error)
	AddEtikettSet(to kennung.EtikettSet, from kennung.EtikettSet) (err error)
	Add(s kennung.EtikettSet) (err error)
}

type Index interface {
	schnittstellen.Flusher
	schnittstellen.Resetter

	EtikettIndex

	hinweis_index.HinweisIndex
}

type index struct {
	path string
	schnittstellen.VerzeichnisseFactory
	etiketten  map[kennung.Etikett]int64
	didRead    bool
	hasChanges bool

	hinweisIndex hinweis_index.HinweisIndex
}

type row struct {
	kennung.Etikett
	count int64
}

func MakeIndex(
	k schnittstellen.Konfig,
	s schnittstellen.Standort,
	vf schnittstellen.VerzeichnisseFactory,
) (i *index, err error) {
	i = &index{
		path:                 s.FileVerzeichnisseEtiketten(),
		VerzeichnisseFactory: vf,
		etiketten:            make(map[kennung.Etikett]int64),
	}

	if i.hinweisIndex, err = hinweis_index.MakeIndex(
		k,
		s,
		vf,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) Flush() (err error) {
	if err = i.hinweisIndex.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !i.hasChanges {
		errors.Log().Print("no changes")
		return
	}

	var w1 io.WriteCloser

	if w1, err = i.WriteCloserVerzeichnisse(i.path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, w1.Close)

	w := bufio.NewWriter(w1)

	defer errors.Deferred(&err, w.Flush)

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

func (i *index) readIfNecessary() (err error) {
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

func (i *index) AddEtikettSet(
	to kennung.EtikettSet,
	from kennung.EtikettSet,
) (err error) {
	d := kennung.MakeSetEtikettDelta(
		from,
		to,
	)

	if err = i.processDelta(d); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) processDelta(d kennung.EtikettDelta) (err error) {
	if err = i.Add(d.Added); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = i.del(d.Removed); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) Add(s kennung.EtikettSet) (err error) {
	if s.Len() == 0 {
		errors.Log().Print("no etiketten to add")
		return
	}

	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.hasChanges = true

	for _, e := range s.Elements() {
		errors.Log().Printf("adding etiketten: %s", e)
		var c int64

		c, _ = i.etiketten[e]
		c += 1
		i.etiketten[e] = c
	}

	return
}

func (i *index) del(s kennung.EtikettSet) (err error) {
	if s.Len() == 0 {
		errors.Log().Print("no etiketten to delete")
		return
	}

	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.hasChanges = true

	for _, e := range s.Elements() {
		errors.Log().Printf("removing etikett: %s", e)
		var c int64
		ok := false

		if c, ok = i.etiketten[e]; !ok {
			errors.Log().Print(errors.Errorf("attempting to delete etikett that is already at 0"))
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

func (i *index) GetAllEtiketten() (es []kennung.Etikett, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	es = make([]kennung.Etikett, len(i.etiketten))

	n := 0

	for e := range i.etiketten {
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

func (i *index) Reset() (err error) {
	if err = i.hinweisIndex.Reset(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) AddHinweis(h kennung.Hinweis) (err error) {
	if err = i.hinweisIndex.AddHinweis(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) CreateHinweis() (h kennung.Hinweis, err error) {
	if h, err = i.hinweisIndex.CreateHinweis(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *index) PeekHinweisen(n int) (hs []kennung.Hinweis, err error) {
	if hs, err = i.hinweisIndex.PeekHinweisen(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
