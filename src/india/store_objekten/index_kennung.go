package store_objekten

import (
	"bufio"
	"encoding/gob"
	"io"

	"github.com/friedenberg/zit/src/alfa/kennung"
	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/golf/hinweisen"
)

type encodedKennung struct {
	AvailableKennung map[int]bool
	Yin, Yang        []string
}

type indexKennung struct {
	*umwelt.Umwelt
	ioFactory

	path string

	encodedKennung

	oldHinweisenStore hinweisen.Hinweisen

	didRead    bool
	hasChanges bool
}

func newIndexKennung(
	u *umwelt.Umwelt,
	ioFactory ioFactory,
	oldHinweisenStore hinweisen.Hinweisen,
	p string,
) (i *indexKennung, err error) {
	i = &indexKennung{
		Umwelt:            u,
		path:              p,
		oldHinweisenStore: oldHinweisenStore,
		ioFactory:         ioFactory,
		encodedKennung: encodedKennung{
			AvailableKennung: make(map[int]bool),
			Yin:              make([]string, 0),
			Yang:             make([]string, 0),
		},
	}

	return
}

func (i *indexKennung) Flush() (err error) {
	if !i.hasChanges {
		logz.Print("no changes")
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

	return
}

func (i *indexKennung) readIfNecessary() (err error) {
	if i.didRead {
		return
	}

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
	logz.Printf("deleting %d, %s", n, h)

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

	var n int

	for n, _ = range i.AvailableKennung {
		break
	}

	if _, ok := i.AvailableKennung[int(n)]; ok {
		delete(i.AvailableKennung, int(n))
	}

	i.hasChanges = true

	k := &kennung.Kennung{}
	k.SetInt(kennung.Int(n))

	return hinweis.New(
		k.Id(),
		i.oldHinweisenStore.Factory().Left(),
		i.oldHinweisenStore.Factory().Right(),
	)
}
