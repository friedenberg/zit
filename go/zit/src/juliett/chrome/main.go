package chrome

import (
	"net/url"
	"sync"

	"code.linenisgreat.com/chrest/go/chrest"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/juliett/konfig"
)

type transacted struct {
	sync.Mutex
	schnittstellen.MutableSetLike[*kennung.Kennung2]
}

type Store struct {
	konfig            *konfig.Compiled
	externalStoreInfo sku.ExternalStoreInfo
	typ               kennung.Typ
	chrome            chrest.Browser

	tabCache cache

	urls map[url.URL][]item

	l       sync.Mutex
	removed map[url.URL]struct{}
	added   map[url.URL][]*kennung.Kennung2

	transacted transacted

	transactedUrlIndex   map[url.URL]sku.TransactedMutableSet
	transactedTabIdIndex map[float64]*sku.Transacted

	itemDeletedStringFormatWriter schnittstellen.FuncIter[*CheckedOut]
}

func MakeChrome(
	k *konfig.Compiled,
	s standort.Standort,
	itemDeletedStringFormatWriter schnittstellen.FuncIter[*CheckedOut],
) *Store {
	c := &Store{
		konfig:  k,
		typ:     kennung.MustTyp("toml-bookmark"),
		removed: make(map[url.URL]struct{}),
		added:   make(map[url.URL][]*kennung.Kennung2),
		transacted: transacted{
			MutableSetLike: collections_value.MakeMutableValueSet(
				iter.StringerKeyer[*kennung.Kennung2]{},
			),
		},
		transactedUrlIndex:            make(map[url.URL]sku.TransactedMutableSet),
		transactedTabIdIndex:          make(map[float64]*sku.Transacted),
		itemDeletedStringFormatWriter: itemDeletedStringFormatWriter,
	}

	return c
}

func (c *Store) GetVirtualStore() sku.ExternalStoreLike {
	return c
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

	if r, err = c.externalStoreInfo.AkteReader(sk.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, r)

	var tb sku_fmt.TomlBookmark

	dec := toml.NewDecoder(r)

	if err = dec.Decode(&tb); err != nil {
		err = errors.Wrapf(err, "Sha: %s, Kennung: %s", sk.GetAkteSha(), sk.GetKennung())
		return
	}

	if u, err = url.Parse(tb.Url); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Store) QueryCheckedOut(
	qg sku.ExternalQuery,
	f schnittstellen.FuncIter[sku.CheckedOutLike],
) (err error) {
	// o := sku.ObjekteOptions{
	// 	Mode: objekte_mode.ModeRealizeSansProto,
	// }

	var co CheckedOut

	for u, items := range c.urls {
		matchingUrls, ok := c.transactedUrlIndex[u]

		for _, item := range items {
			var uChrome *url.URL

			if uChrome, err = item.GetUrl(); err != nil {
				err = errors.Wrap(err)
				return
			}

			processOne := func(internal *sku.Transacted) (err error) {
				co.External.browser.Kennung.SetGattung(gattung.Zettel)
				co.External.Kennung.SetGattung(gattung.Zettel)

				if internal != nil {
					var uSku *url.URL

					if uSku, err = c.getUrl(internal); err != nil {
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

				} else {
					sku.TransactedResetter.Reset(&co.External.Transacted)
					sku.TransactedResetter.Reset(&co.Internal)
					co.State = checked_out_state.StateUntracked
				}

				browser := &co.External.browser
				co.External.item = item

				if co.External.Metadatei.Tai, err = item.GetTai(); err != nil {
					err = errors.Wrap(err)
					return
				}

				browser.Metadatei.Tai = co.External.Metadatei.Tai

				if browser.Metadatei.Typ, err = item.GetTyp(); err != nil {
					err = errors.Wrap(err)
					return
				}

				if browser.Metadatei.Bezeichnung, err = item.GetBezeichnung(); err != nil {
					err = errors.Wrap(err)
					return
				}

				if !qg.ContainsSku(browser) {
					return
				}

				if err = f(&co); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			}

			tabId, okTabId := item.GetTabId()
			var matchingTabId *sku.Transacted

			if okTabId {
				matchingTabId, okTabId = c.transactedTabIdIndex[tabId]
			}

			if !ok || okTabId {
				if err = processOne(matchingTabId); err != nil {
					err = errors.Wrapf(err, "Item: %#v", item)
					return
				}
			} else if ok && !qg.ExcludeUntracked {
				if err = matchingUrls.Each(processOne); err != nil {
					err = errors.Wrap(err)
					return
				}
			}
		}
	}

	return
}

func (c *Store) CheckoutOne(
	options checkout_options.Options,
	sz *sku.Transacted,
) (cz sku.CheckedOutLike, err error) {
	if !sz.Metadatei.Typ.Equals(c.typ) {
		err = errors.Wrap(sku.ErrExternalStoreUnsupportedTyp(sz.Metadatei.Typ))
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
	co.External.browser.Metadatei.Typ = kennung.MustTyp("!chrome-tab")
	co.External.item = map[string]interface{}{"url": u.String()}

	c.l.Lock()
	defer c.l.Unlock()

	existing := c.added[*u]
	c.added[*u] = append(existing, sz.Kennung.Clone())

	// 	ui.Debug().Print(response)

	return
}
