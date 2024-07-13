package store

import (
	"bufio"
	"encoding/gob"
	"io"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/tridex"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO-P4 make generic
type AbbrStore interface {
	Hinweis() AbbrStoreGeneric[kennung.Hinweis, *kennung.Hinweis]
	Kisten() AbbrStoreGeneric[kennung.Kasten, *kennung.Kasten]
	Shas() AbbrStoreGeneric[sha.Sha, *sha.Sha]
	Etiketten() AbbrStoreGeneric[kennung.Etikett, *kennung.Etikett]
	Typen() AbbrStoreGeneric[kennung.Typ, *kennung.Typ]

	AddMatchable(*sku.Transacted) error
	GetAbbr() kennung.Abbr

	errors.Flusher
}

type indexAbbrEncodableTridexes struct {
	Shas      indexNotHinweis[sha.Sha, *sha.Sha]
	Hinweis   indexHinweis
	Etiketten indexNotHinweis[kennung.Etikett, *kennung.Etikett]
	Typen     indexNotHinweis[kennung.Typ, *kennung.Typ]
	Kisten    indexNotHinweis[kennung.Kasten, *kennung.Kasten]
}

type indexAbbr struct {
	lock     sync.Locker
	once     *sync.Once
	standort standort.Standort

	path string

	indexAbbrEncodableTridexes

	didRead    bool
	hasChanges bool
}

func newIndexAbbr(
	standort standort.Standort,
	p string,
) (i *indexAbbr, err error) {
	i = &indexAbbr{
		lock:     &sync.Mutex{},
		once:     &sync.Once{},
		path:     p,
		standort: standort,
		indexAbbrEncodableTridexes: indexAbbrEncodableTridexes{
			Shas: indexNotHinweis[sha.Sha, *sha.Sha]{
				Kennungen: tridex.Make(),
			},
			Hinweis: indexHinweis{
				Kopfen:    tridex.Make(),
				Schwanzen: tridex.Make(),
			},
			Etiketten: indexNotHinweis[kennung.Etikett, *kennung.Etikett]{
				Kennungen: tridex.Make(),
			},
			Typen: indexNotHinweis[kennung.Typ, *kennung.Typ]{
				Kennungen: tridex.Make(),
			},
			Kisten: indexNotHinweis[kennung.Kasten, *kennung.Kasten]{
				Kennungen: tridex.Make(),
			},
		},
	}

	i.indexAbbrEncodableTridexes.Hinweis.readFunc = i.readIfNecessary
	i.indexAbbrEncodableTridexes.Kisten.readFunc = i.readIfNecessary
	i.indexAbbrEncodableTridexes.Shas.readFunc = i.readIfNecessary
	i.indexAbbrEncodableTridexes.Etiketten.readFunc = i.readIfNecessary
	i.indexAbbrEncodableTridexes.Typen.readFunc = i.readIfNecessary

	return
}

func (i *indexAbbr) Flush() (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if !i.hasChanges {
		ui.Log().Print("no changes")
		return
	}

	var w1 io.WriteCloser

	if w1, err = i.standort.WriteCloserCache(i.path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w1)

	w := bufio.NewWriter(w1)

	defer errors.DeferredFlusher(&err, w)

	enc := gob.NewEncoder(w)

	if err = enc.Encode(i.indexAbbrEncodableTridexes); err != nil {
		err = errors.Wrapf(err, "failed to write encoded kennung")
		return
	}

	return
}

func (i *indexAbbr) readIfNecessary() (err error) {
	i.once.Do(
		func() {
			if i.didRead {
				return
			}

			ui.Log().Print("reading")

			i.didRead = true

			var r1 io.ReadCloser

			if r1, err = i.standort.ReadCloserCache(i.path); err != nil {
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

			ui.Log().Print("starting decode")

			if err = dec.Decode(&i.indexAbbrEncodableTridexes); err != nil {
				ui.Log().Print("finished decode unsuccessfully")
				err = errors.Wrap(err)
				return
			}
		},
	)

	return
}

func (i *indexAbbr) GetAbbr() (out kennung.Abbr) {
	out.Hinweis.Expand = i.Hinweis().ExpandStringString
	out.Sha.Expand = i.Shas().ExpandStringString

	out.Hinweis.Abbreviate = i.Hinweis().Abbreviate
	out.Sha.Abbreviate = i.Shas().Abbreviate

	return
}

func (i *indexAbbr) AddMatchable(o *sku.Transacted) (err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.hasChanges = true

	i.indexAbbrEncodableTridexes.Shas.Kennungen.Add(o.GetAkteSha().String())

	ks := o.GetKennung().String()

	switch o.GetGenre() {
	case gattung.Zettel:
		var h kennung.Hinweis

		if err = h.Set(ks); err != nil {
			err = errors.Wrap(err)
			return
		}

		i.indexAbbrEncodableTridexes.Hinweis.Kopfen.Add(h.GetHead())
		i.indexAbbrEncodableTridexes.Hinweis.Schwanzen.Add(h.GetTail())

	case gattung.Typ:
		i.indexAbbrEncodableTridexes.Typen.Kennungen.Add(ks)

	case gattung.Etikett:
		i.indexAbbrEncodableTridexes.Etiketten.Kennungen.Add(ks)

	case gattung.Kasten:
		i.indexAbbrEncodableTridexes.Kisten.Kennungen.Add(ks)

		// default:
		// 	err = errors.Errorf("unsupported objekte: %T", to)
		// 	return
	}

	return
}

// TODO switch to using ennui for existence
func (i *indexAbbr) Exists(k *kennung.Kennung2) (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	switch k.GetGenre() {
	case gattung.Zettel:
		err = i.Hinweis().Exists(k.Parts())

	case gattung.Typ:
		err = i.Typen().Exists(k.Parts())

	case gattung.Etikett:
		err = i.Etiketten().Exists(k.Parts())

	case gattung.Kasten:
		err = i.Kisten().Exists(k.Parts())

	case gattung.Konfig:
		// konfig always exists
		return

	default:
		err = collections.MakeErrNotFound(k)
	}

	return
}

func (i *indexAbbr) Hinweis() (asg AbbrStoreGeneric[kennung.Hinweis, *kennung.Hinweis]) {
	asg = &i.indexAbbrEncodableTridexes.Hinweis

	return
}

func (i *indexAbbr) Kisten() (asg AbbrStoreGeneric[kennung.Kasten, *kennung.Kasten]) {
	asg = &i.indexAbbrEncodableTridexes.Kisten

	return
}

func (i *indexAbbr) Shas() (asg AbbrStoreGeneric[sha.Sha, *sha.Sha]) {
	asg = &i.indexAbbrEncodableTridexes.Shas

	return
}

func (i *indexAbbr) Etiketten() (asg AbbrStoreGeneric[kennung.Etikett, *kennung.Etikett]) {
	asg = &i.indexAbbrEncodableTridexes.Etiketten

	return
}

func (i *indexAbbr) Typen() (asg AbbrStoreGeneric[kennung.Typ, *kennung.Typ]) {
	asg = &i.indexAbbrEncodableTridexes.Typen

	return
}
