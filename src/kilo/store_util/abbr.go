package store_util

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/tridex"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type AbbrStore interface {
	HinweisExists(kennung.Hinweis) error
	ExpandShaString(string) (sha.Sha, error)
	ExpandEtikettString(string) (kennung.Etikett, error)
	ExpandHinweisString(string) (kennung.Hinweis, error)
	AbbreviateSha(schnittstellen.Value) (string, error)
	AbbreviateHinweis(schnittstellen.Korper) (string, error)
	AddStoredAbbreviation(schnittstellen.Stored) error
	errors.Flusher
}

type indexAbbrEncodableTridexes struct {
	Shas             *tridex.Tridex
	HinweisKopfen    *tridex.Tridex
	HinweisSchwanzen *tridex.Tridex
	Etiketten        *tridex.Tridex
}

type indexAbbr struct {
	lock sync.Locker
	StoreUtilVerzeichnisse

	path string

	indexAbbrEncodableTridexes

	didRead    bool
	hasChanges bool
}

func newIndexAbbr(
	suv StoreUtilVerzeichnisse,
	p string,
) (i AbbrStore, err error) {
	i = &indexAbbr{
		lock:                   &sync.Mutex{},
		path:                   p,
		StoreUtilVerzeichnisse: suv,
		indexAbbrEncodableTridexes: indexAbbrEncodableTridexes{
			Shas:             tridex.Make(),
			HinweisKopfen:    tridex.Make(),
			HinweisSchwanzen: tridex.Make(),
			Etiketten:        tridex.Make(),
		},
	}

	return
}

func (i *indexAbbr) Flush() (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

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

	if err = enc.Encode(i.indexAbbrEncodableTridexes); err != nil {
		err = errors.Wrapf(err, "failed to write encoded kennung")
		return
	}

	return
}

func (i *indexAbbr) readIfNecessary() (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if i.didRead {
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

	defer errors.Deferred(&err, r1.Close)

	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	errors.Log().Print("starting decode")

	if err = dec.Decode(&i.indexAbbrEncodableTridexes); err != nil {
		errors.Log().Print("finished decode unsuccessfully")
		err = errors.Wrap(err)
		return
	}

	errors.Log().Print("finished decode successfully")

	return
}

func (i *indexAbbr) AddStoredAbbreviation(o schnittstellen.Stored) (err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.hasChanges = true

	i.indexAbbrEncodableTridexes.Shas.Add(o.GetAkteSha().String())
	i.indexAbbrEncodableTridexes.Shas.Add(o.GetObjekteSha().String())

	if z, ok := o.(*zettel.Transacted); ok {
		i.indexAbbrEncodableTridexes.HinweisKopfen.Add(z.Kennung().Kopf())
		i.indexAbbrEncodableTridexes.HinweisSchwanzen.Add(z.Kennung().Schwanz())

		for _, e := range kennung.Expanded(z.Objekte.Etiketten, kennung.ExpanderRight).Elements() {
			i.indexAbbrEncodableTridexes.Etiketten.Add(e.String())
		}
	}

	return
}

func (i *indexAbbr) AbbreviateSha(s schnittstellen.Value) (abbr string, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	abbr = s.String()

	if i.GetKonfig().PrintAbbreviatedShas {
		abbr = i.indexAbbrEncodableTridexes.Shas.Abbreviate(abbr)
	}

	return
}

func (i *indexAbbr) ExpandShaString(st string) (s sha.Sha, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	expanded := i.indexAbbrEncodableTridexes.Shas.Expand(st)

	if err = s.Set(expanded); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *indexAbbr) HinweisExists(h kennung.Hinweis) (err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !i.indexAbbrEncodableTridexes.HinweisKopfen.ContainsExactly(h.Kopf()) {
		err = objekte_store.ErrNotFound{Id: h}
		return
	}

	if !i.indexAbbrEncodableTridexes.HinweisSchwanzen.ContainsExactly(h.Schwanz()) {
		err = objekte_store.ErrNotFound{Id: h}
		return
	}

	return
}

func (i *indexAbbr) AbbreviateHinweis(h schnittstellen.Korper) (v string, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	kopf := h.Kopf()
	schwanz := h.Schwanz()

	if i.GetKonfig().PrintAbbreviatedHinweisen {
		kopf = i.indexAbbrEncodableTridexes.HinweisKopfen.Abbreviate(h.Kopf())
		schwanz = i.indexAbbrEncodableTridexes.HinweisSchwanzen.Abbreviate(h.Schwanz())

		if kopf == "" || schwanz == "" {
			err = errors.Errorf("abbreviated kopf would be empty for %s", h)
			errors.Log().PrintDebug(i.indexAbbrEncodableTridexes.HinweisKopfen)
			return
		}

		if schwanz == "" {
			err = errors.Errorf("abbreviated schwanz would be empty for %s", h)
			return
		}
	}

	v = fmt.Sprintf("%s/%s", kopf, schwanz)

	return
}

func (i *indexAbbr) ExpandHinweisString(s string) (h kennung.Hinweis, err error) {
	errors.Log().Print(s)

	var ha kennung.Hinweis

	if ha, err = kennung.MakeHinweis(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	return i.ExpandHinweis(ha)
}

func (i *indexAbbr) ExpandHinweis(hAbbr kennung.Hinweis) (h kennung.Hinweis, err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	kopf := i.indexAbbrEncodableTridexes.HinweisKopfen.Expand(hAbbr.Kopf())
	schwanz := i.indexAbbrEncodableTridexes.HinweisSchwanzen.Expand(hAbbr.Schwanz())

	if h, err = kennung.MakeHinweisKopfUndSchwanz(kopf, schwanz); err != nil {
		err = errors.Wrapf(err, "{Abbreviated: '%s'}", hAbbr)
		return
	}

	return
}

func (i *indexAbbr) ExpandEtikettString(s string) (e kennung.Etikett, err error) {
	errors.Log().Print(s)

	if e, err = kennung.MakeEtikett(s); err != nil {
		err = errors.Wrap(err)
		return
	}

	return i.ExpandEtikett(e)
}

func (i *indexAbbr) ExpandEtikett(eAbbr kennung.Etikett) (e kennung.Etikett, err error) {
	errors.Log().Print(eAbbr)

	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ex := i.indexAbbrEncodableTridexes.Etiketten.Expand(eAbbr.String())

	if ex == "" {
		//TODO should try to use the expansion if possible
		ex = eAbbr.String()
	}

	if err = e.Set(ex); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
