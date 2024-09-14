package store_browser

import (
	"net/url"
	"sync"

	"code.linenisgreat.com/chrest/go/src/charlie/browser_items"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/juliett/config"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

type transacted struct {
	sync.Mutex
	interfaces.MutableSetLike[*ids.ObjectId]
}

type Store struct {
	config            *config.Compiled
	externalStoreInfo external_store.Info
	typ               ids.Type
	browser           browser_items.BrowserProxy

	tabCache cache

	urls map[url.URL][]Item

	l       sync.Mutex
	deleted map[url.URL][]Item
	added   map[url.URL][]Item

	itemsById map[string]Item

	transacted transacted

	transactedUrlIndex  map[url.URL]sku.TransactedMutableSet
	transactedItemIndex map[browser_items.ItemId]*sku.Transacted

	itemDeletedStringFormatWriter interfaces.FuncIter[Item]
}

func Make(
	k *config.Compiled,
	s fs_home.Home,
	itemDeletedStringFormatWriter interfaces.FuncIter[Item],
) *Store {
	c := &Store{
		config:    k,
		typ:       ids.MustType("toml-bookmark"),
		deleted:   make(map[url.URL][]Item),
		added:     make(map[url.URL][]Item),
		itemsById: make(map[string]Item),
		transacted: transacted{
			MutableSetLike: collections_value.MakeMutableValueSet(
				iter.StringerKeyer[*ids.ObjectId]{},
			),
		},
		transactedUrlIndex:            make(map[url.URL]sku.TransactedMutableSet),
		transactedItemIndex:           make(map[browser_items.ItemId]*sku.Transacted),
		itemDeletedStringFormatWriter: itemDeletedStringFormatWriter,
	}

	return c
}

func (fs *Store) GetExternalStoreLike() external_store.StoreLike {
	return fs
}

func (s *Store) ApplyDotOperator() error {
	return nil
}

func (s *Store) GetObjectIdsForString(
	v string,
) (k []sku.ExternalObjectId, err error) {
	item, ok := s.itemsById[v]

	if !ok {
		err = errors.Errorf("not a browser item id")
		return
	}

	k = append(k, item.GetExternalObjectId())

	return
}

func (s *Store) Flush() (err error) {
	wg := iter.MakeErrorWaitGroupParallel()

	wg.Do(s.flushUrls)

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Store) getUrl(sk *sku.Transacted) (u *url.URL, err error) {
	var r sha.ReadCloser

	if r, err = c.externalStoreInfo.BlobReader(sk.GetBlobSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, r)

	var tb sku_fmt.TomlBookmark

	dec := toml.NewDecoder(r)

	if err = dec.Decode(&tb); err != nil {
		err = errors.Wrapf(err, "Sha: %s, Object Id: %s", sk.GetBlobSha(), sk.GetObjectId())
		return
	}

	if u, err = url.Parse(tb.Url); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Store) CheckoutOne(
	options checkout_options.Options,
	sz *sku.Transacted,
) (cz sku.CheckedOutLike, err error) {
	if !sz.Metadata.Type.Equals(c.typ) {
		err = errors.Wrap(external_store.ErrUnsupportedTyp(sz.Metadata.Type))
		return
	}

	var u *url.URL

	if u, err = c.getUrl(sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	co := GetCheckedOutPool().Get()
	cz = co
	var item Item
	item.Url.URL = *u
	item.ExternalId = sz.ObjectId.String()
	item.Id.Type = "tab"

	sku.TransactedResetter.ResetWith(co.GetSku(), sz)
	sku.TransactedResetter.ResetWith(co.GetSkuExternalLike().GetSku(), sz)
	co.State = checked_out_state.JustCheckedOut
	co.External.ExternalType = ids.MustType("!browser-tab")

	if err = item.WriteToExternal(&co.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.l.Lock()
	defer c.l.Unlock()

	existing := c.added[*u]
	c.added[*u] = append(existing, item)

	// 	ui.Debug().Print(response)

	return
}

func (c *Store) QueryCheckedOut(
	qg *query.Group,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	// o := sku.ObjekteOptions{
	// 	Mode: objekte_mode.ModeRealizeSansProto,
	// }

	ex := executor{
		store: c,
		qg:    qg,
		out:   f,
	}

	for u, items := range c.urls {
		matchingUrls, exactIndexURLMatch := c.transactedUrlIndex[u]

		for _, item := range items {
			var matchingTabId *sku.Transacted
			var trackedFromBefore bool

			tabId := item.Id
			matchingTabId, trackedFromBefore = c.transactedItemIndex[tabId]

			if trackedFromBefore {
				if err = ex.tryToEmitOneExplicitlyCheckedOut(
					matchingTabId,
					item,
				); err != nil {
					err = errors.Wrapf(err, "Item: %#v", item)
					return
				}
			} else if !exactIndexURLMatch {
				if err = ex.tryToEmitOneUntracked(item); err != nil {
					err = errors.Wrapf(err, "Item: %#v", item)
					return
				}
			} else if exactIndexURLMatch {
				if err = matchingUrls.Each(
					func(matching *sku.Transacted) (err error) {
						if err = ex.tryToEmitOneRecognized(
							matching,
							item,
						); err != nil {
							err = errors.Wrapf(err, "Item: %#v", item)
							return
						}

						return
					},
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}
		}
	}

	return
}

func (c *Store) GetExternalStoreOrganizeFormat(
	f *sku_fmt.Box,
) sku_fmt.ExternalLike {
	fo := MakeFormatOrganize(f)

	return sku_fmt.ExternalLike{
		ReaderExternalLike: fo,
		WriterExternalLike: fo,
	}
}

func (c *Store) GetExternalLikePool() interfaces.PoolValue[sku.ExternalLike] {
	return pool.Bespoke[sku.ExternalLike]{
		FuncGet: func() sku.ExternalLike {
			return poolExternal.Get()
		},
		FuncPut: func(e sku.ExternalLike) {
			poolExternal.Put(e.(*External))
		},
	}
}

func (c *Store) GetExternalLikeResetter3() interfaces.Resetter3[sku.ExternalLike] {
	return pool.BespokeResetter[sku.ExternalLike]{
		FuncReset: func(el sku.ExternalLike) {
			a := el.(*External)
			sku.ExternalResetter.Reset(&a.External)
		},
		FuncResetWith: func(eldst, elsrc sku.ExternalLike) {
			dst, src := eldst.(*External), elsrc.(*External)
			sku.ExternalResetter.ResetWith(&dst.External, &src.External)
		},
	}
}
