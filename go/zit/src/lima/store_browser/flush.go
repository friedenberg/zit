package store_browser

import (
	"syscall"

	"code.linenisgreat.com/chrest/go/src/charlie/browser_items"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"golang.org/x/exp/maps"
)

func (s *Store) flushUrls() (err error) {
	if len(s.deleted) == 0 && len(s.added) == 0 {
		return
	}

	var resp browser_items.HTTPResponseWithRequestPayloadPut

	deleted := make(map[string]checkedOutWithItem, len(s.deleted))

	var req browser_items.BrowserRequestPut
	req.Deleted = make([]browser_items.Item, 0, len(s.deleted))

	for _, is := range s.deleted {
		for _, i := range is {
			req.Deleted = append(req.Deleted, i.Item.Item)
			deleted[i.Item.Item.ExternalId] = i
		}
	}

	for _, is := range s.added {
		for _, i := range is {
			req.Added = append(req.Added, i.Item.Item)
		}
	}

	if !s.config.DryRun {
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
	} else {
		for _, is := range s.deleted {
			for _, i := range is {
				resp.Deleted = append(resp.Deleted, i.Item.Item)
			}
		}

		for _, is := range s.added {
			for _, i := range is {
				resp.Added = append(resp.Added, i.Item.Item)
			}
		}
	}

	for _, i := range resp.RequestPayloadPut.Added {
		// TODO emit changes
		s.tabCache.Rows[i.ExternalId] = i.Id
	}

	for _, item := range resp.RequestPayloadPut.Deleted {
		delete(s.tabCache.Rows, item.ExternalId)

		originalItem, ok := deleted[item.ExternalId]

		if !ok {
			err = errors.Errorf(
				"missing item with id %q from deleted cache: %q",
				item.ExternalId,
				maps.Keys(deleted),
			)

			return
		}

		if s.config.IsDryRun() {
			originalItem.CheckedOut.GetSkuExternal().State = external_state.WouldDelete
		} else {
			originalItem.CheckedOut.GetSkuExternal().State = external_state.Deleted
		}

		if err = s.itemDeletedStringFormatWriter(
			originalItem.CheckedOut,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	clear(s.added)
	clear(s.deleted)

	if err = s.flushCache(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
