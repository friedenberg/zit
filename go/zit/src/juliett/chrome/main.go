package chrome

import (
	"net/url"
	"sync"

	"code.linenisgreat.com/chrest/go/chrest"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

type transacted struct {
	sync.Mutex
	interfaces.MutableSetLike[*ids.ObjectId]
}

type Store struct {
	konfig            *konfig.Compiled
	externalStoreInfo external_store.Info
	typ               ids.Type
	chrome            chrest.Browser

	tabCache cache

	urls map[url.URL][]item

	l       sync.Mutex
	removed map[url.URL]struct{}
	added   map[url.URL][]*ids.ObjectId

	transacted transacted

	transactedUrlIndex   map[url.URL]sku.TransactedMutableSet
	transactedTabIdIndex map[float64]*sku.Transacted

	itemDeletedStringFormatWriter interfaces.FuncIter[*CheckedOut]
}

func MakeChrome(
	k *konfig.Compiled,
	s fs_home.Home,
	itemDeletedStringFormatWriter interfaces.FuncIter[*CheckedOut],
) *Store {
	c := &Store{
		konfig:  k,
		typ:     ids.MustType("toml-bookmark"),
		removed: make(map[url.URL]struct{}),
		added:   make(map[url.URL][]*ids.ObjectId),
		transacted: transacted{
			MutableSetLike: collections_value.MakeMutableValueSet(
				iter.StringerKeyer[*ids.ObjectId]{},
			),
		},
		transactedUrlIndex:            make(map[url.URL]sku.TransactedMutableSet),
		transactedTabIdIndex:          make(map[float64]*sku.Transacted),
		itemDeletedStringFormatWriter: itemDeletedStringFormatWriter,
	}

	return c
}

func (fs *Store) GetExternalStoreLike() external_store.StoreLike {
	return fs
}

func (c *Store) GetExternalKennung() (ks interfaces.SetLike[*ids.ObjectId], err error) {
	ksm := collections_value.MakeMutableValueSet[*ids.ObjectId](nil)
	ks = ksm

	for u, items := range c.urls {
		for _, item := range items {
			var matchingTabId *sku.Transacted
			var trackedFromBefore bool

			{
				tabId, okTabId := item.GetTabId()

				if okTabId {
					matchingTabId, trackedFromBefore = c.transactedTabIdIndex[tabId]
				}
			}

			if trackedFromBefore {
				if err = ksm.Add(matchingTabId.Kennung.Clone()); err != nil {
					err = errors.Wrap(err)
					return
				}
			} else {
				k := ids.GetObjectIdPool().Get()

				if err = k.SetRaw(u.String()); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = ksm.Add(k); err != nil {
					err = errors.Wrap(err)
					return
				}
			}
		}
	}

	return
}

// TODO
func (s *Store) GetKennungForString(v string) (k *ids.ObjectId, err error) {
	err = collections.MakeErrNotFoundString(v)
	return
	k = ids.GetObjectIdPool().Get()

	if err = k.SetRaw(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) Flush() (err error) {
	if s.konfig.DryRun {
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

	if r, err = c.externalStoreInfo.BlobReader(sk.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, r)

	var tb sku_fmt.TomlBookmark

	dec := toml.NewDecoder(r)

	if err = dec.Decode(&tb); err != nil {
		err = errors.Wrapf(err, "Sha: %s, Kennung: %s", sk.GetAkteSha(), sk.GetObjectId())
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
	if !sz.Metadatei.Type.Equals(c.typ) {
		err = errors.Wrap(sku.ErrExternalStoreUnsupportedTyp(sz.Metadatei.Type))
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
	co.State = checked_out_state.StateJustCheckedOut
	co.External.browser.Metadatei.Type = ids.MustType("!chrome-tab")
	co.External.item = map[string]interface{}{"url": u.String()}

	c.l.Lock()
	defer c.l.Unlock()

	existing := c.added[*u]
	c.added[*u] = append(existing, sz.Kennung.Clone())

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

	var co CheckedOut

	for u, items := range c.urls {
		matchingUrls, exactIndexURLMatch := c.transactedUrlIndex[u]

		for _, item := range items {
			var matchingTabId *sku.Transacted
			var trackedFromBefore bool

			{
				tabId, okTabId := item.GetTabId()

				if okTabId {
					matchingTabId, trackedFromBefore = c.transactedTabIdIndex[tabId]
				}
			}

			if trackedFromBefore {
				if err = c.tryToEmitOneExplicitlyCheckedOut(
					qg,
					matchingTabId,
					&co,
					item,
					f,
				); err != nil {
					err = errors.Wrapf(err, "Item: %#v", item)
					return
				}
			} else if !exactIndexURLMatch {
				if err = c.tryToEmitOneUntracked(
					qg,
					&co,
					item,
					f,
				); err != nil {
					err = errors.Wrapf(err, "Item: %#v", item)
					return
				}
			} else if exactIndexURLMatch {
				if err = matchingUrls.Each(
					func(matching *sku.Transacted) (err error) {
						if err = c.tryToEmitOneRecognized(
							qg,
							matching,
							&co,
							item,
							f,
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

func (c *Store) QueryUnsure(
	qg *query.Group,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	// o := sku.ObjekteOptions{
	// 	Mode: objekte_mode.ModeRealizeSansProto,
	// }

	var co CheckedOut

	for u, items := range c.urls {
		matchingUrls, exactIndexURLMatch := c.transactedUrlIndex[u]

		for _, item := range items {
			{
				tabId, okTabId := item.GetTabId()

				if okTabId {
					if _, trackedFromBefore := c.transactedTabIdIndex[tabId]; trackedFromBefore {
						continue
					}
				}
			}

			if !exactIndexURLMatch {
				if err = c.tryToEmitOneUntracked(
					qg,
					&co,
					item,
					f,
				); err != nil {
					err = errors.Wrapf(err, "Item: %#v", item)
					return
				}
			} else if exactIndexURLMatch {
				if err = matchingUrls.Each(
					func(matching *sku.Transacted) (err error) {
						if err = c.tryToEmitOneRecognized(
							qg,
							matching,
							&co,
							item,
							f,
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

func (c *Store) tryToEmitOneExplicitlyCheckedOut(
	qg *query.Group,
	internal *sku.Transacted,
	co *CheckedOut,
	item item,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	sku.TransactedResetter.Reset(&co.External.browser)
	co.External.browser.Kennung.SetGenre(genres.Zettel)
	co.External.Kennung.SetGenre(genres.Zettel)

	var uSku *url.URL

	if uSku, err = c.getUrl(internal); err != nil {
		err = errors.Wrap(err)
		return
	}

	var uChrome *url.URL

	if uChrome, err = item.GetUrl(); err != nil {
		err = errors.Wrap(err)
		return
	}

	sku.TransactedResetter.ResetWith(&co.Internal, internal)
	sku.TransactedResetter.ResetWith(&co.External.Transacted, internal)

	if *uSku == *uChrome {
		co.State = checked_out_state.StateExistsAndSame
	} else {
		co.State = checked_out_state.StateExistsAndDifferent
	}

	if err = c.tryToEmitOneCommon(
		qg,
		co,
		item,
		false,
		f,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Store) tryToEmitOneRecognized(
	qg *query.Group,
	internal *sku.Transacted,
	co *CheckedOut,
	item item,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	if !qg.IncludeRecognized {
		return
	}

	sku.TransactedResetter.Reset(&co.External.browser)
	co.External.browser.Kennung.SetGenre(genres.Zettel)
	co.External.Kennung.SetGenre(genres.Zettel)

	sku.TransactedResetter.ResetWith(&co.Internal, internal)
	sku.TransactedResetter.ResetWith(&co.External.Transacted, internal)

	co.State = checked_out_state.StateRecognized

	if err = c.tryToEmitOneCommon(
		qg,
		co,
		item,
		true,
		f,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Store) tryToEmitOneUntracked(
	qg *query.Group,
	co *CheckedOut,
	item item,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	if qg.ExcludeUntracked {
		return
	}

	sku.TransactedResetter.Reset(&co.External.browser)
	co.External.browser.Kennung.SetGenre(genres.Zettel)
	co.External.Kennung.SetGenre(genres.Zettel)

	sku.TransactedResetter.Reset(&co.External.Transacted)
	sku.TransactedResetter.Reset(&co.Internal)
	co.State = checked_out_state.StateUntracked

	if err = c.tryToEmitOneCommon(
		qg,
		co,
		item,
		true,
		f,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Store) tryToEmitOneCommon(
	qg *query.Group,
	co *CheckedOut,
	i item,
	overwrite bool,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	browser := &co.External.browser

	if err = co.External.SetItem(i, overwrite); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !qg.ContainsSku(browser) && !qg.ContainsSku(co.GetSku()) {
		return
	}

	if err = f(co); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
