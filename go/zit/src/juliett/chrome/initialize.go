package chrome

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"syscall"

	"code.linenisgreat.com/chrest/go/chrest"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

func (s *Store) Initialize(esi external_store.Info) (err error) {
	s.externalStoreInfo = esi

	if err = s.chrome.Read(); err != nil {
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
	var resp chrest.ResponseWithParsedJSONBody
	var req chrest.BrowserRequest

	req.Method = "GET"
	req.Path = "/urls"

	if resp, err = s.request(req); err != nil {
		if errors.IsErrno(err, syscall.ECONNREFUSED) {
			if !s.konfig.Quiet {
				ui.Err().Print("chrest offline")
			}

			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	var chromeTabsRaw2 []interface{}

	switch t := resp.ParsedJSONBody.(type) {
	case []interface{}:
		chromeTabsRaw2 = t

	// case nil:
	// 	return

	default:
		err = errors.Errorf(
			"expected %T, but got %T, %#v",
			chromeTabsRaw2,
			resp.ParsedJSONBody,
			resp.ParsedJSONBody,
		)

		return
	}

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

	s.urls = chromeTabs

	return
}

func (s *Store) flushUrls() (err error) {
	if len(s.removed) == 0 && len(s.added) == 0 {
		return
	}

	b := bytes.NewBuffer(nil)

	req := chrest.BrowserRequest{
		Method: "PUT",
		Path:   "/urls",
		Body:   io.NopCloser(b),
	}

	var reqPayload requestUrlsPut
	reqPayload.Deleted = make([]string, 0, len(s.removed))

	for u := range s.removed {
		reqPayload.Deleted = append(reqPayload.Deleted, u.String())
	}

	lookup := make([][]*kennung.Id, 0, len(s.added))

	for u, k := range s.added {
		reqPayload.Added = append(
			reqPayload.Added,
			createOneTabRequest{
				Url: u.String(),
			},
		)

		lookup = append(lookup, k)
	}

	enc := json.NewEncoder(b)

	if err = enc.Encode(reqPayload); err != nil {
		err = errors.Wrap(err)
		return
	}

	req.Body = io.NopCloser(b)

	var resp chrest.ResponseWithParsedJSONBody

	if resp, err = s.request(req); err != nil {
		if errors.IsErrno(err, syscall.ECONNREFUSED) {
			ui.Err().Print("chrest offline")
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	// TODO get req header for launch time and compare with our lauch time
	// parse response and add urls to cache
	if resp.ParsedJSONBody == nil {
		err = errors.Errorf("got nil response")
		return
	}

	json, ok := resp.ParsedJSONBody.(chrest.JSONObject)

	if !ok {
		err = errors.Errorf(
			"expected %T but got %T, %#v",
			json,
			resp.ParsedJSONBody,
			resp.ParsedJSONBody,
		)

		return
	}

	added, ok := json["added"].(chrest.JSONArray)

	if !ok {
		err = errors.Errorf(
			"expected %T but got %T, %#v",
			added,
			json["added"],
			json["added"],
		)

		return
	}

	if len(added) != len(lookup) {
		err = errors.Errorf("expected to create %d tabs, but got %d tabs", len(lookup), len(added))
		return
	}

	for i, t := range lookup {
		a := added[i].(chrest.JSONObject)

		for _, k := range t {
			s.tabCache.Rows[k.String()] = a["id"].(float64)
		}
	}

	ui.Debug().Printf("%#v", added)

	clear(s.added)
	clear(s.removed)

	if err = s.flushCache(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
