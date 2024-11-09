package store

import (
	"bufio"
	"encoding/gob"
	"io"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/charlie/tridex"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO-P4 make generic to just ObjectIds
type AbbrStore interface {
	ZettelId() AbbrStoreGeneric[ids.ZettelId, *ids.ZettelId]
	Shas() AbbrStoreGeneric[sha.Sha, *sha.Sha]

	AddObjectToAbbreviationStore(*sku.Transacted) error
	GetAbbr() ids.Abbr

	errors.Flusher
}

type indexAbbrEncodableTridexes struct {
	Shas     indexNotZettelId[sha.Sha, *sha.Sha]
	ZettelId indexZettelId
}

type indexAbbr struct {
	options_print.V0

	lock      sync.Locker
	once      *sync.Once
	dirLayout dir_layout.DirLayout

	path string

	indexAbbrEncodableTridexes

	didRead    bool
	hasChanges bool
}

func newIndexAbbr(
	options options_print.V0,
	dirLayout dir_layout.DirLayout,
	p string,
) (i *indexAbbr, err error) {
	i = &indexAbbr{
		V0:        options,
		lock:      &sync.Mutex{},
		once:      &sync.Once{},
		path:      p,
		dirLayout: dirLayout,
		indexAbbrEncodableTridexes: indexAbbrEncodableTridexes{
			Shas: indexNotZettelId[sha.Sha, *sha.Sha]{
				ObjectIds: tridex.Make(),
			},
			ZettelId: indexZettelId{
				Kopfen:    tridex.Make(),
				Schwanzen: tridex.Make(),
			},
		},
	}

	i.indexAbbrEncodableTridexes.ZettelId.readFunc = i.readIfNecessary
	i.indexAbbrEncodableTridexes.Shas.readFunc = i.readIfNecessary

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

	if w1, err = i.dirLayout.WriteCloserCache(i.path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w1)

	w := bufio.NewWriter(w1)

	defer errors.DeferredFlusher(&err, w)

	enc := gob.NewEncoder(w)

	if err = enc.Encode(i.indexAbbrEncodableTridexes); err != nil {
		err = errors.Wrapf(err, "failed to write encoded object id")
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

			if r1, err = i.dirLayout.ReadCloserCache(i.path); err != nil {
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

func (i *indexAbbr) GetAbbr() (out ids.Abbr) {
	out.ZettelId.Expand = i.ZettelId().ExpandStringString
	out.Sha.Expand = i.Shas().ExpandStringString

	if i.Abbreviations.ZettelIds {
		out.ZettelId.Abbreviate = i.ZettelId().Abbreviate
	}

	if i.Abbreviations.Shas {
		out.Sha.Abbreviate = i.Shas().Abbreviate
	}

	return
}

func (i *indexAbbr) AddObjectToAbbreviationStore(o *sku.Transacted) (err error) {
	if err = i.readIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.hasChanges = true

	i.indexAbbrEncodableTridexes.Shas.ObjectIds.Add(o.GetBlobSha().String())

	ks := o.GetObjectId().String()

	switch o.GetGenre() {
	case genres.Zettel:
		var h ids.ZettelId

		if err = h.SetFromIdParts(o.GetObjectId().Parts()); err != nil {
			err = errors.Wrap(err)
			return
		}

		i.indexAbbrEncodableTridexes.ZettelId.Kopfen.Add(h.GetHead())
		i.indexAbbrEncodableTridexes.ZettelId.Schwanzen.Add(h.GetTail())

	case genres.Type, genres.Tag, genres.Config, genres.InventoryList:
		return

	default:
		err = errors.Errorf("unsupported object id: %#v", ks)
		return
	}

	return
}

func (i *indexAbbr) ZettelId() (asg AbbrStoreGeneric[ids.ZettelId, *ids.ZettelId]) {
	asg = &i.indexAbbrEncodableTridexes.ZettelId

	return
}

func (i *indexAbbr) Shas() (asg AbbrStoreGeneric[sha.Sha, *sha.Sha]) {
	asg = &i.indexAbbrEncodableTridexes.Shas

	return
}
