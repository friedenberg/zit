package store_browser

import (
	"net/url"
	"sync"

	"code.linenisgreat.com/chrest/go/chrest/src/charlie/browser_items"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
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

	urls map[url.URL][]browserItem

	l       sync.Mutex
	removed map[url.URL][]browserItem
	added   map[url.URL][]browserItem

	itemsById map[string]browserItem

	transacted transacted

	transactedUrlIndex  map[url.URL]sku.TransactedMutableSet
	transactedItemIndex map[browser_items.ItemId]*sku.Transacted

	itemDeletedStringFormatWriter interfaces.FuncIter[browserItem]
}

func Make(
	k *config.Compiled,
	s fs_home.Home,
	itemDeletedStringFormatWriter interfaces.FuncIter[browserItem],
) *Store {
	c := &Store{
		config:    k,
		typ:       ids.MustType("toml-bookmark"),
		removed:   make(map[url.URL][]browserItem),
		added:     make(map[url.URL][]browserItem),
		itemsById: make(map[string]browserItem),
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

func (s *Store) GetObjectIdsForString(v string) (k []sku.ExternalObjectId, err error) {
	item, ok := s.itemsById[v]

	if !ok {
		err = errors.Errorf("not a browser item id")
		return
	}

	k = append(k, &item)

	return
}

func (s *Store) Flush() (err error) {
	if s.config.DryRun {
		return
	}

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

	sku.TransactedResetter.ResetWith(co.GetSku(), sz)
	sku.TransactedResetter.ResetWith(co.GetSkuExternalLike().GetSku(), sz)
	co.State = checked_out_state.JustCheckedOut
	co.External.browser.Metadata.Type = ids.MustType("!browser-tab")
	co.External.browserItem.Url = u.String()
	co.External.browserItem.ExternalId = sz.ObjectId.String()
	co.External.browserItem.Id.Type = "tab"

	c.l.Lock()
	defer c.l.Unlock()

	existing := c.added[*u]
	c.added[*u] = append(existing, co.External.browserItem)

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
	f *sku_fmt.Organize,
) sku_fmt.ExternalLike {
	fo := MakeFormatOrganize(f)

	return sku_fmt.ExternalLike{
		ReaderExternalLike: fo,
		WriterExternalLike: fo,
	}
}