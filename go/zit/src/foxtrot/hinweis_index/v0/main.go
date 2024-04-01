package hinweis_index

import (
	"bufio"
	"encoding/gob"
	"io"
	"math/rand"
	"sync"

	"code.linenisgreat.com/zit/src/alfa/coordinates"
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/delta/hinweisen"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

type encodedKennung struct {
	AvailableKennung map[int]bool
}

type oldIndex struct {
	su schnittstellen.VerzeichnisseFactory

	lock sync.Mutex
	path string

	encodedKennung

	oldHinweisenStore *hinweisen.Hinweisen

	didRead    bool
	hasChanges bool

	nonRandomSelection bool
}

func MakeIndex(
	k schnittstellen.Konfig,
	s schnittstellen.Standort,
	su schnittstellen.VerzeichnisseFactory,
) (i *oldIndex, err error) {
	i = &oldIndex{
		path:               s.FileVerzeichnisseKennung(),
		nonRandomSelection: k.UsePredictableHinweisen(),
		su:                 su,
		encodedKennung: encodedKennung{
			AvailableKennung: make(map[int]bool),
		},
	}

	if i.oldHinweisenStore, err = hinweisen.New(s); err != nil {
		if errors.IsNotExist(err) {
			errors.TodoP4("determine which layer handles no-create kasten")
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (i *oldIndex) Flush() (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if !i.hasChanges {
		errors.Log().Print("no changes")
		return
	}

	var w1 io.WriteCloser

	if w1, err = i.su.WriteCloserVerzeichnisse(i.path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w1)

	w := bufio.NewWriter(w1)

	defer errors.DeferredFlusher(&err, w)

	enc := gob.NewEncoder(w)

	if err = enc.Encode(i.encodedKennung); err != nil {
		err = errors.Wrapf(err, "failed to write encoded kennung")
		return
	}

	return
}

func (i *oldIndex) readIfNecessary() (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if i.didRead {
		return
	}

	errors.Log().Print("reading")

	i.didRead = true

	var r1 io.ReadCloser

	if r1, err = i.su.ReadCloserVerzeichnisse(i.path); err != nil {
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

	if err = dec.Decode(&i.encodedKennung); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *oldIndex) Reset() (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	lMax := i.oldHinweisenStore.Left().Len() - 1
	rMax := i.oldHinweisenStore.Right().Len() - 1

	if lMax == 0 {
		err = errors.Errorf("left kennung are empty")
		return
	}

	if rMax == 0 {
		err = errors.Errorf("right kennung are empty")
		return
	}

	i.AvailableKennung = make(map[int]bool, lMax*rMax)

	for l := 0; l <= lMax; l++ {
		for r := 0; r <= rMax; r++ {
			k := &coordinates.Kennung{
				Left:  coordinates.Int(l),
				Right: coordinates.Int(r),
			}

			errors.Log().Print(k)

			n := int(k.Id())
			i.AvailableKennung[n] = true
		}
	}

	i.hasChanges = true

	return
}

func (i *oldIndex) AddHinweis(k1 kennung.Kennung) (err error) {
	if !k1.GetGattung().EqualsGattung(gattung.Zettel) {
		err = gattung.MakeErrUnsupportedGattung(k1)
		return
	}

	var h kennung.Hinweis

	if err = h.Set(k1.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var left, right int

	if left, err = i.oldHinweisenStore.Left().Kennung(h.Kopf()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if right, err = i.oldHinweisenStore.Right().Kennung(h.Schwanz()); err != nil {
		err = errors.Wrap(err)
		return
	}

	k := coordinates.Kennung{
		Left:  coordinates.Int(left),
		Right: coordinates.Int(right),
	}

	n := k.Id()
	errors.Log().Printf("deleting %d, %s", n, h)

	i.lock.Lock()
	defer i.lock.Unlock()

	delete(i.AvailableKennung, int(n))

	i.hasChanges = true

	return
}

func (i *oldIndex) CreateHinweis() (h *kennung.Hinweis, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if len(i.AvailableKennung) == 0 {
		err = errors.Errorf("no available kennung")
		return
	}

	if len(i.AvailableKennung) == 0 {
		err = errors.Wrap(hinweisen.ErrHinweisenExhausted{})
		return
	}

	ri := 0

	if len(i.AvailableKennung) > 1 {
		ri = rand.Intn(len(i.AvailableKennung) - 1)
	}

	m := 0
	j := 0

	for n := range i.AvailableKennung {
		if i.nonRandomSelection {
			if m == 0 {
				m = n
				continue
			}

			if n > m {
				continue
			}

			m = n
		} else {
			j++
			m = n

			if j == ri {
				break
			}
		}
	}

	delete(i.AvailableKennung, int(m))

	i.hasChanges = true

	return i.makeHinweisButDontStore(m)
}

func (i *oldIndex) makeHinweisButDontStore(
	j int,
) (h *kennung.Hinweis, err error) {
	k := &coordinates.Kennung{}
	k.SetInt(coordinates.Int(j))

	h, err = kennung.NewHinweis(
		k.Id(),
		i.oldHinweisenStore.Left(),
		i.oldHinweisenStore.Right(),
	)

	if err != nil {
		err = errors.Wrapf(err, "trying to make hinweis for %s", k)
		return
	}

	return
}

func (i *oldIndex) PeekHinweisen(m int) (hs []*kennung.Hinweis, err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if m > len(i.AvailableKennung) || m == 0 {
		m = len(i.AvailableKennung)
	}

	hs = make([]*kennung.Hinweis, 0, m)
	j := 0

	for n := range i.AvailableKennung {
		k := &coordinates.Kennung{}
		k.SetInt(coordinates.Int(n))

		var h *kennung.Hinweis

		if h, err = i.makeHinweisButDontStore(n); err != nil {
			err = errors.Wrap(err)
			return
		}

		hs = append(hs, h)

		j++

		if j == m {
			break
		}
	}

	return
}
