package store_objekten

import (
	"bufio"
	"encoding/gob"
	"io"
	"math/rand"
	"time"

	"github.com/friedenberg/zit/src/alfa/coordinates"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/hinweisen"
	"github.com/friedenberg/zit/src/delta/konfig"
)

type encodedKennung struct {
	AvailableKennung map[int]bool
}

type indexKennung struct {
	ioFactory

	path string

	encodedKennung

	oldHinweisenStore *hinweisen.Hinweisen

	didRead    bool
	hasChanges bool

	nonRandomSelection bool
}

func newIndexKennung(
	k konfig.Konfig,
	ioFactory ioFactory,
	oldHinweisenStore *hinweisen.Hinweisen,
	p string,
) (i *indexKennung, err error) {
	i = &indexKennung{
		path:               p,
		nonRandomSelection: k.PredictableHinweisen,
		oldHinweisenStore:  oldHinweisenStore,
		ioFactory:          ioFactory,
		encodedKennung: encodedKennung{
			AvailableKennung: make(map[int]bool),
		},
	}

	return
}

func (i *indexKennung) Flush() (err error) {
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

	if err = enc.Encode(i.encodedKennung); err != nil {
		err = errors.Wrapf(err, "failed to write encoded kennung")
		return
	}

	return
}

func (i *indexKennung) readIfNecessary() (err error) {
	if i.didRead {
		errors.Log().Print("already read")
		return
	}

	errors.Log().Print("reading")

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

	if err = dec.Decode(&i.encodedKennung); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *indexKennung) reset() (err error) {
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

func (i *indexKennung) addHinweis(h hinweis.Hinweis) (err error) {
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

	if _, ok := i.AvailableKennung[int(n)]; ok {
		delete(i.AvailableKennung, int(n))
	}

	i.hasChanges = true

	return
}

func (i *indexKennung) createHinweis() (h hinweis.Hinweis, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if len(i.AvailableKennung) == 0 {
		err = errors.Errorf("no available kennung")
		return
	}

	rand.Seed(time.Now().UnixNano())

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

	for n, _ := range i.AvailableKennung {
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

	if _, ok := i.AvailableKennung[int(m)]; ok {
		delete(i.AvailableKennung, int(m))
	}

	i.hasChanges = true

	return i.makeHinweisButDontStore(m)
}

func (i *indexKennung) makeHinweisButDontStore(j int) (h hinweis.Hinweis, err error) {
	k := &coordinates.Kennung{}
	k.SetInt(coordinates.Int(j))

	h, err = hinweis.New(
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

func (i *indexKennung) PeekHinweisen(m int) (hs []hinweis.Hinweis, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if m > len(i.AvailableKennung) || m == 0 {
		m = len(i.AvailableKennung)
	}

	hs = make([]hinweis.Hinweis, 0, m)
	j := 0

	for n, _ := range i.AvailableKennung {
		k := &coordinates.Kennung{}
		k.SetInt(coordinates.Int(n))

		var h hinweis.Hinweis

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
