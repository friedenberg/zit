package kennung

import (
	"bufio"
	"encoding/gob"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/friedenberg/zit/src/alfa/coordinates"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/hinweisen"
)

type hinweisIndex struct {
	su schnittstellen.VerzeichnisseFactory

	lock *sync.RWMutex
	path string

	bitset collections.Bitset

	oldHinweisenStore *hinweisen.Hinweisen

	didRead    bool
	hasChanges bool

	nonRandomSelection bool
}

func makeHinweisIndex(
	k schnittstellen.Konfig,
	s Standort,
	su schnittstellen.VerzeichnisseFactory,
) (i *hinweisIndex, err error) {
	i = &hinweisIndex{
		lock:               &sync.RWMutex{},
		path:               s.FileVerzeichnisseHinweis(),
		nonRandomSelection: k.UsePredictableHinweisen(),
		su:                 su,
		bitset:             collections.MakeBitset(0),
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

func (i *hinweisIndex) Flush() (err error) {
	i.lock.RLock()

	if !i.hasChanges {
		errors.Log().Print("no changes")
		i.lock.RUnlock()
		return
	}

	i.lock.RUnlock()

	var w1 io.WriteCloser

	if w1, err = i.su.WriteCloserVerzeichnisse(i.path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, w1.Close)

	w := bufio.NewWriter(w1)

	defer errors.Deferred(&err, w.Flush)

	enc := gob.NewEncoder(w)

	if err = enc.Encode(i.bitset); err != nil {
		err = errors.Wrapf(err, "failed to write encoded kennung")
		return
	}

	return
}

func (i *hinweisIndex) readIfNecessary() (err error) {
	i.lock.RLock()

	if i.didRead {
		i.lock.RUnlock()
		return
	}

	i.lock.RUnlock()

	i.lock.Lock()
	defer i.lock.Unlock()

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

	if err = dec.Decode(i.bitset); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *hinweisIndex) Reset() (err error) {
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

	i.bitset = collections.MakeBitsetOn(lMax * rMax)

	i.hasChanges = true

	return
}

func (i *hinweisIndex) AddHinweis(h Hinweis) (err error) {
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

	i.bitset.DelIfPresent(int(n))

	i.hasChanges = true

	return
}

func (i *hinweisIndex) CreateHinweis() (h Hinweis, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if i.bitset.CountOn() == 0 {
		err = errors.Errorf("no available kennung")
		return
	}

	rand.Seed(time.Now().UnixNano())

	if i.bitset.CountOn() == 0 {
		err = errors.Wrap(hinweisen.ErrHinweisenExhausted{})
		return
	}

	ri := 0

	if i.bitset.CountOn() > 1 {
		ri = rand.Intn(i.bitset.CountOn() - 1)
	}

	m := 0
	j := 0

	if err = i.bitset.Each(
		func(n int) (err error) {
			if i.nonRandomSelection {
				if m == 0 {
					m = n
					return
				}

				if n > m {
					return
				}

				m = n
			} else {
				j++
				m = n

				if j == ri {
					err = collections.ErrStopIteration
					return
				}
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.bitset.DelIfPresent(int(m))

	i.hasChanges = true

	return i.makeHinweisButDontStore(m)
}

func (i *hinweisIndex) makeHinweisButDontStore(j int) (h Hinweis, err error) {
	k := &coordinates.Kennung{}
	k.SetInt(coordinates.Int(j))

	h, err = NewHinweis(
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

func (i *hinweisIndex) PeekHinweisen(m int) (hs []Hinweis, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if m > i.bitset.CountOn() || m == 0 {
		m = i.bitset.CountOn()
	}

	hs = make([]Hinweis, 0, m)
	j := 0

	if err = i.bitset.Each(
		func(n int) (err error) {
			k := &coordinates.Kennung{}
			k.SetInt(coordinates.Int(n))

			var h Hinweis

			if h, err = i.makeHinweisButDontStore(n); err != nil {
				err = errors.Wrap(err)
				return
			}

			hs = append(hs, h)

			j++

			if j == m {
				err = collections.ErrStopIteration
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
