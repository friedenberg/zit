package chrome

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"syscall"
	"time"

	"code.linenisgreat.com/chrest/go/chrest"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (c *Store) Initialize() (err error) {
	if !c.konfig.ChrestEnabled {
		return
	}

	wg := iter.MakeErrorWaitGroupParallel()

	wg.Do(c.initializeUrls)
	wg.Do(c.initializeIndex)

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO abstract and regenerate on commit / reindex
func (c *Store) initializeIndex() (err error) {
	var l sync.Mutex

	if err = c.storeFuncs.FuncQuery(
		nil,
		func(sk *sku.Transacted) (err error) {
			if !sk.GetTyp().Equals(c.typ) {
				return
			}

			var u *url.URL

			if u, err = c.getUrl(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			cl := sku.GetTransactedPool().Get()
			sku.TransactedResetter.ResetWith(cl, sk)

			l.Lock()
			defer l.Unlock()

			existing, ok := c.transactedUrlIndex[*u]

			if !ok {
				existing = sku.MakeTransactedMutableSet()
				c.transactedUrlIndex[*u] = existing
			}

			if err = existing.Add(cl); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Store) initializeUrls() (err error) {
	if err = c.chrestConfig.Read(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var chromeTabsRaw interface{}
	var req *http.Request

	if req, err = http.NewRequest("GET", "http://localhost/urls", nil); err != nil {
		err = errors.Wrap(err)
		return
	}

	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Duration(1e9)),
	)

	defer cancel()

	if chromeTabsRaw, err = chrest.AskChrome(ctx, c.chrestConfig, req); err != nil {
		if errors.IsErrno(err, syscall.ECONNREFUSED) {
			if !c.konfig.Quiet {
				ui.Err().Print("chrest offline")
			}

			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	chromeTabsRaw2 := chromeTabsRaw.([]interface{})

	chromeTabs := make(map[url.URL][]item, len(chromeTabsRaw2))

	for _, tabRaw := range chromeTabsRaw2 {
		tab := tabRaw.(map[string]interface{})
		ur := tab["url"]

		if ur == nil {
			continue
		}

		var u *url.URL

		if u, err = url.Parse(ur.(string)); err != nil {
			err = errors.Wrap(err)
			return
		}

		chromeTabs[*u] = append(chromeTabs[*u], tab)
	}

	c.urls = chromeTabs

	return
}
