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
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO-P4 make generic
type AbbrStore interface {
	ZettelId() AbbrStoreGeneric[ids.ZettelId, *ids.ZettelId]
	Kisten() AbbrStoreGeneric[ids.RepoId, *ids.RepoId]
	Shas() AbbrStoreGeneric[sha.Sha, *sha.Sha]
	Etiketten() AbbrStoreGeneric[ids.Tag, *ids.Tag]
	Typen() AbbrStoreGeneric[ids.Type, *ids.Type]

	AddMatchable(*sku.Transacted) error
	GetAbbr() ids.Abbr

	errors.Flusher
}

type indexAbbrEncodableTridexes struct {
	Shas     indexNotZettelId[sha.Sha, *sha.Sha]
	ZettelId indexZettelId
	Tags     indexNotZettelId[ids.Tag, *ids.Tag]
	Types    indexNotZettelId[ids.Type, *ids.Type]
	RepoIds  indexNotZettelId[ids.RepoId, *ids.RepoId]
}

type indexAbbr struct {
	lock    sync.Locker
	once    *sync.Once
	fs_home fs_home.Home

	path string

	indexAbbrEncodableTridexes

	didRead    bool
	hasChanges bool
}

func newIndexAbbr(
	fs_home fs_home.Home,
	p string,
) (i *indexAbbr, err error) {
	i = &indexAbbr{
		lock:    &sync.Mutex{},
		once:    &sync.Once{},
		path:    p,
		fs_home: fs_home,
		indexAbbrEncodableTridexes: indexAbbrEncodableTridexes{
			Shas: indexNotZettelId[sha.Sha, *sha.Sha]{
				ObjectIds: tridex.Make(),
			},
			ZettelId: indexZettelId{
				Kopfen:    tridex.Make(),
				Schwanzen: tridex.Make(),
			},
			Tags: indexNotZettelId[ids.Tag, *ids.Tag]{
				ObjectIds: tridex.Make(),
			},
			Types: indexNotZettelId[ids.Type, *ids.Type]{
				ObjectIds: tridex.Make(),
			},
			RepoIds: indexNotZettelId[ids.RepoId, *ids.RepoId]{
				ObjectIds: tridex.Make(),
			},
		},
	}

	i.indexAbbrEncodableTridexes.ZettelId.readFunc = i.readIfNecessary
	i.indexAbbrEncodableTridexes.RepoIds.readFunc = i.readIfNecessary
	i.indexAbbrEncodableTridexes.Shas.readFunc = i.readIfNecessary
	i.indexAbbrEncodableTridexes.Tags.readFunc = i.readIfNecessary
	i.indexAbbrEncodableTridexes.Types.readFunc = i.readIfNecessary

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

	if w1, err = i.fs_home.WriteCloserCache(i.path); err != nil {
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

			if r1, err = i.fs_home.ReadCloserCache(i.path); err != nil {
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

	out.ZettelId.Abbreviate = i.ZettelId().Abbreviate
	out.Sha.Abbreviate = i.Shas().Abbreviate

	return
}

func (i *indexAbbr) AddMatchable(o *sku.Transacted) (err error) {
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

		if err = h.Set(ks); err != nil {
			err = errors.Wrap(err)
			return
		}

		i.indexAbbrEncodableTridexes.ZettelId.Kopfen.Add(h.GetHead())
		i.indexAbbrEncodableTridexes.ZettelId.Schwanzen.Add(h.GetTail())

	case genres.Type:
		i.indexAbbrEncodableTridexes.Types.ObjectIds.Add(ks)

	case genres.Tag:
		i.indexAbbrEncodableTridexes.Tags.ObjectIds.Add(ks)

	case genres.Repo:
		i.indexAbbrEncodableTridexes.RepoIds.ObjectIds.Add(ks)

		// default:
		// 	err = errors.Errorf("unsupported objekte: %T", to)
		// 	return
	}

	return
}

// TODO switch to using ennui for existence
func (i *indexAbbr) Exists(k *ids.ObjectId) (err error) {
	i.lock.Lock()
	defer i.lock.Unlock()

	switch k.GetGenre() {
	case genres.Zettel:
		err = i.ZettelId().Exists(k.Parts())

	case genres.Type:
		err = i.Typen().Exists(k.Parts())

	case genres.Tag:
		err = i.Etiketten().Exists(k.Parts())

	case genres.Repo:
		err = i.Kisten().Exists(k.Parts())

	case genres.Config:
		// konfig always exists
		return

	default:
		err = collections.MakeErrNotFound(k)
	}

	return
}

func (i *indexAbbr) ZettelId() (asg AbbrStoreGeneric[ids.ZettelId, *ids.ZettelId]) {
	asg = &i.indexAbbrEncodableTridexes.ZettelId

	return
}

func (i *indexAbbr) Kisten() (asg AbbrStoreGeneric[ids.RepoId, *ids.RepoId]) {
	asg = &i.indexAbbrEncodableTridexes.RepoIds

	return
}

func (i *indexAbbr) Shas() (asg AbbrStoreGeneric[sha.Sha, *sha.Sha]) {
	asg = &i.indexAbbrEncodableTridexes.Shas

	return
}

func (i *indexAbbr) Etiketten() (asg AbbrStoreGeneric[ids.Tag, *ids.Tag]) {
	asg = &i.indexAbbrEncodableTridexes.Tags

	return
}

func (i *indexAbbr) Typen() (asg AbbrStoreGeneric[ids.Type, *ids.Type]) {
	asg = &i.indexAbbrEncodableTridexes.Types

	return
}
