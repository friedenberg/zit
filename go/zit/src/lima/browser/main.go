package browser

import (
	"net/url"
	"sync"

	"code.linenisgreat.com/chrest/go/chrest"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
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
	browser           chrest.Browser

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

func Make(
	k *config.Compiled,
	s fs_home.Home,
	itemDeletedStringFormatWriter interfaces.FuncIter[*CheckedOut],
) *Store {
	c := &Store{
		config:  k,
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

// TODO
func (s *Store) GetObjectIdsForString(v string) (k []sku.ExternalObjectId, err error) {
	err = errors.Implement()
	return
	// k = []sku.ExternalObjectId{ids.GetObjectIdPool().Get()}

	// if err = k[0].SetRaw(v); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// return
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
	co.External.item = map[string]interface{}{"url": u.String()}

	c.l.Lock()
	defer c.l.Unlock()

	existing := c.added[*u]
	c.added[*u] = append(existing, sz.ObjectId.Clone())

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

func (c *Store) tryToEmitOneExplicitlyCheckedOut(
	qg *query.Group,
	internal *sku.Transacted,
	co *CheckedOut,
	item item,
	f interfaces.FuncIter[sku.CheckedOutLike],
) (err error) {
	sku.TransactedResetter.Reset(&co.External.browser)
	co.External.browser.ObjectId.SetGenre(genres.Zettel)
	co.External.ObjectId.SetGenre(genres.Zettel)

	var uSku *url.URL

	if uSku, err = c.getUrl(internal); err != nil {
		err = errors.Wrap(err)
		return
	}

	var uBrowser *url.URL

	if uBrowser, err = item.GetUrl(); err != nil {
		err = errors.Wrap(err)
		return
	}

	sku.TransactedResetter.ResetWith(&co.Internal, internal)
	sku.TransactedResetter.ResetWith(&co.External.Transacted, internal)

	if *uSku == *uBrowser {
		co.State = checked_out_state.ExistsAndSame
	} else {
		co.State = checked_out_state.ExistsAndDifferent
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
	co.State = checked_out_state.Recognized

	if !qg.ContainsSkuCheckedOutState(co.State) {
		return
	}

	sku.TransactedResetter.Reset(&co.External.browser)
	co.External.browser.ObjectId.SetGenre(genres.Unknown)
	co.External.ObjectId.SetGenre(genres.Unknown)

	if err = item.WriteToObjectId(&co.External.ObjectId); err != nil {
		err = errors.Wrap(err)
		return
	}

	sku.TransactedResetter.ResetWith(&co.Internal, internal)
	sku.TransactedResetter.ResetWith(&co.External.Transacted, internal)

	co.State = checked_out_state.Recognized

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
	co.State = checked_out_state.Untracked

	if !qg.ContainsSkuCheckedOutState(co.State) {
		return
	}

	sku.TransactedResetter.Reset(&co.External.browser)

	if err = item.WriteToObjectId(&co.External.ObjectId); err != nil {
		err = errors.Wrap(err)
		return
	}

	co.External.browser.ObjectId.SetGenre(genres.Zettel)
	co.External.ObjectId.SetGenre(genres.Zettel)

	sku.TransactedResetter.Reset(&co.External.Transacted)
	sku.TransactedResetter.Reset(&co.Internal)

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

	if !qg.ContainsExternalSku(browser, co.State) &&
		!qg.ContainsExternalSku(co.GetSku(), co.State) {
		return
	}

	if err = co.External.ObjectId.SetRepoId(
		c.externalStoreInfo.RepoId.String(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = co.Internal.ObjectId.SetRepoId(
		c.externalStoreInfo.RepoId.String(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = f(co); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Store) GetExternalStoreOrganizeFormat(
	f *sku_fmt.Organize,
) sku_fmt.ExternalLike {
	return MakeFormatOrganize(f)
}
