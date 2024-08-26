package store_browser

import (
	"net/url"
	"syscall"

	"code.linenisgreat.com/chrest/go/chrest/src/charlie/browser_items"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

func (s *Store) Initialize(esi external_store.Info) (err error) {
	s.externalStoreInfo = esi

	if err = s.browser.Read(); err != nil {
		err = errors.Wrap(err)
		return
	}

	wg := iter.MakeErrorWaitGroupParallel()

	wg.Do(s.initializeUrls)
	wg.Do(s.initializeIndex)

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) initializeUrls() (err error) {
	var req browser_items.BrowserRequestGet
	var resp browser_items.HTTPResponseWithRequestPayloadGet

	if resp, err = s.browser.Get(req); err != nil {
		if errors.IsErrno(err, syscall.ECONNREFUSED) {
			if !s.config.Quiet {
				ui.Err().Print("chrest offline")
			}

			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	s.urls = make(map[url.URL][]browserItem, len(resp.RequestPayloadGet))

	if err = s.resetCacheIfNecessary(resp.Response); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, item := range resp.RequestPayloadGet {
		i := browserItem{Item: item}

		var u *url.URL

		if u, err = i.GetUrl(); err != nil {
			err = errors.Wrap(err)
			return
		}

		s.urls[*u] = append(s.urls[*u], i)
		s.itemsById[i.GetObjectId().String()] = i
	}

	return
}

func (s *Store) flushUrls() (err error) {
	if len(s.removed) == 0 && len(s.added) == 0 {
		return
	}

	var req browser_items.BrowserRequestPut
	req.Deleted = make([]browser_items.Item, 0, len(s.removed))

	for _, is := range s.removed {
		for _, i := range is {
			req.Deleted = append(req.Deleted, i.Item)
		}
	}

	for _, is := range s.added {
		for _, i := range is {
			req.Added = append(req.Added, i.Item)
		}
	}

	var resp browser_items.HTTPResponseWithRequestPayloadPut

	if resp, err = s.browser.Put(req); err != nil {
		if errors.IsErrno(err, syscall.ECONNREFUSED) {
			ui.Err().Print("chrest offline")
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.resetCacheIfNecessary(resp.Response); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, i := range resp.RequestPayloadPut.Added {
		s.tabCache.Rows[i.ExternalId] = i.Id
	}

	for _, i := range resp.RequestPayloadPut.Deleted {
		delete(s.tabCache.Rows, i.ExternalId)

		if err = s.itemDeletedStringFormatWriter(browserItem{Item: i}); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	clear(s.added)
	clear(s.removed)

	if err = s.flushCache(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
