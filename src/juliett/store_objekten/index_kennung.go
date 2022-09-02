package store_objekten

import (
	"bufio"
	"encoding/gob"
	"io"
	"math/rand"
	"time"

	"github.com/friedenberg/zit/src/alfa/kennung"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/golf/hinweisen"
)

type encodedKennung struct {
	AvailableKennung map[int]bool
}

type indexKennung struct {
	*umwelt.Umwelt
	ioFactory

	path string

	encodedKennung

	oldHinweisenStore hinweisen.Hinweisen

	didRead    bool
	hasChanges bool

	nonRandomSelection bool
}

func newIndexKennung(
	u *umwelt.Umwelt,
	ioFactory ioFactory,
	oldHinweisenStore hinweisen.Hinweisen,
	p string,
) (i *indexKennung, err error) {
	i = &indexKennung{
		Umwelt:             u,
		path:               p,
		nonRandomSelection: u.Konfig.PredictableHinweisen,
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
		errors.Print("no changes")
		return
	}

	var w1 io.WriteCloser

	if w1, err = i.WriteCloserVerzeichnisse(i.path); err != nil {
		err = errors.Error(err)
		return
	}

	defer stdprinter.PanicIfError(w1.Close)

	w := bufio.NewWriter(w1)

	defer stdprinter.PanicIfError(w.Flush)

	enc := gob.NewEncoder(w)

	if err = enc.Encode(i.encodedKennung); err != nil {
		err = errors.Wrapped(err, "failed to write encoded kennung")
		return
	}

	errors.PrintDebug(i.encodedKennung)

	return
}

func (i *indexKennung) readIfNecessary() (err error) {
	if i.didRead {
		errors.Print("already read")
		return
	}

	errors.Print("reading")

	i.didRead = true

	var r1 io.ReadCloser

	if r1, err = i.ReadCloserVerzeichnisse(i.path); err != nil {
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

	if err = dec.Decode(&i.encodedKennung); err != nil {
		err = errors.Error(err)
		return
	}

	errors.PrintDebug(i.encodedKennung)

	return
}

func (i *indexKennung) reset() (err error) {
	lMax := i.oldHinweisenStore.Factory().Left().Len() - 1
	rMax := i.oldHinweisenStore.Factory().Right().Len() - 1

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
			k := &kennung.Kennung{
				Left:  kennung.Int(l),
				Right: kennung.Int(r),
			}

			errors.Print(k)

			n := int(k.Id())
			i.AvailableKennung[n] = true
		}
	}

	i.hasChanges = true

	return
}

func (i *indexKennung) addHinweis(h hinweis.Hinweis) (err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	var left, right int

	if left, err = i.oldHinweisenStore.Factory().Left().Kennung(h.Kopf()); err != nil {
		err = errors.Error(err)
		return
	}

	if right, err = i.oldHinweisenStore.Factory().Right().Kennung(h.Schwanz()); err != nil {
		err = errors.Error(err)
		return
	}

	k := kennung.Kennung{
		Left:  kennung.Int(left),
		Right: kennung.Int(right),
	}

	n := k.Id()
	errors.Printf("deleting %d, %s", n, h)

	if _, ok := i.AvailableKennung[int(n)]; ok {
		delete(i.AvailableKennung, int(n))
	}

	i.hasChanges = true

	return
}

func (i *indexKennung) createHinweis() (h hinweis.Hinweis, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	if len(i.AvailableKennung) == 0 {
		err = errors.Errorf("no available kennung")
		return
	}

	rand.Seed(time.Now().UnixNano())
	ri := rand.Intn(len(i.AvailableKennung) - 1)

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
	k := &kennung.Kennung{}
	k.SetInt(kennung.Int(j))

	h, err = hinweis.New(
		k.Id(),
		i.oldHinweisenStore.Factory().Left(),
		i.oldHinweisenStore.Factory().Right(),
	)

	if err != nil {
		err = errors.Wrapped(err, "trying to make hinweis for %s", k)
		return
	}

	return
}

func (i *indexKennung) PeekHinweisen(m int) (hs []hinweis.Hinweis, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Error(err)
		return
	}

	if m > len(i.AvailableKennung) || m == 0 {
		m = len(i.AvailableKennung)
	}

	hs = make([]hinweis.Hinweis, 0, m)
	j := 0

	for n, _ := range i.AvailableKennung {
		k := &kennung.Kennung{}
		k.SetInt(kennung.Int(n))

		var h hinweis.Hinweis

		if h, err = i.makeHinweisButDontStore(n); err != nil {
			err = errors.Error(err)
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
