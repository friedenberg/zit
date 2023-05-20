package store_util

import (
	"bufio"
	"encoding/gob"
	"io"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/tridex"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/kasten"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

// TODO-P4 make generic
type AbbrStore interface {
	Hinweis() AbbrStoreGeneric[kennung.Hinweis]
	Kisten() AbbrStoreGeneric[kennung.Kasten]
	Shas() AbbrStoreGeneric[sha.Sha]
	Etiketten() AbbrStoreGeneric[kennung.Etikett]
	Typen() AbbrStoreGeneric[kennung.Typ]

	AddMatchable(kennung.Matchable) error

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
	lock sync.Locker
	once *sync.Once
	StoreUtilVerzeichnisse

	path string

	indexAbbrEncodableTridexes

	didRead    bool
	hasChanges bool
}

func newIndexAbbr(
	suv StoreUtilVerzeichnisse,
	p string,
) (i *indexAbbr, err error) {
	i = &indexAbbr{
		lock:                   &sync.Mutex{},
		once:                   &sync.Once{},
		path:                   p,
		StoreUtilVerzeichnisse: suv,
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
	i.once.Do(
		func() {
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
		},
	)

	return
}

func (i *indexAbbr) AddMatchable(o kennung.Matchable) (err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.hasChanges = true

	i.indexAbbrEncodableTridexes.Shas.Kennungen.Add(o.GetAkteSha().String())
	i.indexAbbrEncodableTridexes.Shas.Kennungen.Add(o.GetObjekteSha().String())

	switch to := o.(type) {
	case *zettel.Transacted:
		i.indexAbbrEncodableTridexes.Hinweis.Kopfen.Add(to.Kennung().Kopf())
		i.indexAbbrEncodableTridexes.Hinweis.Schwanzen.Add(to.Kennung().Schwanz())

	case *typ.Transacted:
		i.indexAbbrEncodableTridexes.Typen.Add(to.Sku.Kennung)

	case *etikett.Transacted:
		i.indexAbbrEncodableTridexes.Etiketten.Kennungen.Add(to.Sku.Kennung.String())

	case *kasten.Transacted:
		i.indexAbbrEncodableTridexes.Kisten.Kennungen.Add(to.Sku.Kennung.String())

		// default:
		// 	err = errors.Errorf("unsupported objekte: %T", to)
		// 	return
	}

	return
}

func (i *indexAbbr) Hinweis() (asg AbbrStoreGeneric[kennung.Hinweis]) {
	asg = &i.indexAbbrEncodableTridexes.Hinweis

	if !i.GetKonfig().PrintAbbreviatedHinweisen {
		asg = indexNoAbbr[kennung.Hinweis, *kennung.Hinweis]{
			AbbrStoreGeneric: asg,
		}
	}

	return
}

func (i *indexAbbr) Kisten() (asg AbbrStoreGeneric[kennung.Kasten]) {
	asg = &i.indexAbbrEncodableTridexes.Kisten

	if !i.GetKonfig().PrintAbbreviatedKennungen {
		asg = indexNoAbbr[kennung.Kasten, *kennung.Kasten]{
			AbbrStoreGeneric: asg,
		}
	}

	return
}

func (i *indexAbbr) Shas() (asg AbbrStoreGeneric[sha.Sha]) {
	asg = &i.indexAbbrEncodableTridexes.Shas

	if !i.GetKonfig().PrintAbbreviatedShas {
		asg = indexNoAbbr[sha.Sha, *sha.Sha]{
			AbbrStoreGeneric: asg,
		}
	}

	return
}

func (i *indexAbbr) Etiketten() (asg AbbrStoreGeneric[kennung.Etikett]) {
	asg = &i.indexAbbrEncodableTridexes.Etiketten

	if !i.GetKonfig().PrintAbbreviatedKennungen {
		asg = indexNoAbbr[kennung.Etikett, *kennung.Etikett]{
			AbbrStoreGeneric: asg,
		}
	}

	return
}

func (i *indexAbbr) Typen() (asg AbbrStoreGeneric[kennung.Typ]) {
	asg = &i.indexAbbrEncodableTridexes.Typen

	if !i.GetKonfig().PrintAbbreviatedKennungen {
		asg = indexNoAbbr[kennung.Typ, *kennung.Typ]{
			AbbrStoreGeneric: asg,
		}
	}

	return
}
