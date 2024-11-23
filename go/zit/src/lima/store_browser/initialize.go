package store_browser

import (
	"context"
	"net/url"
	"syscall"

	"code.linenisgreat.com/chrest/go/src/charlie/browser_items"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

func (s *Store) Initialize(esi external_store.Supplies) (err error) {
	s.externalStoreInfo = esi

	if err = s.browser.Read(); err != nil {
		err = errors.Wrap(err)
		return
	}

	wg := quiter.MakeErrorWaitGroupParallel()

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

	ui.Log().Print("getting all")

  ctx := context.Background()
  ctxWithTimeout, cancel := context.WithTimeout(ctx, 1e9)
  defer cancel()

	if resp, err = s.browser.GetAll(
    ctxWithTimeout,
    req,
  ); err != nil {
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

	ui.Log().Print("got all")

	s.urls = make(map[url.URL][]Item, len(resp.RequestPayloadGet))

	if err = s.resetCacheIfNecessary(resp.Response); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, item := range resp.RequestPayloadGet {
		i := Item{Item: item}

		u := i.Url.URL

		s.urls[u] = append(s.urls[u], i)
		s.itemsById[i.GetObjectId().String()] = i
	}

	return
}
