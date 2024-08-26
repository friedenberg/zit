package store_browser

import (
	"fmt"
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type browserItemRaw map[string]interface{}

func (ct browserItemRaw) GetTagSet() ids.TagSet {
	me := ids.MakeTagMutableSet()

	switch ct["type"].(string) {
	case "history":
		me.Add(
			ids.MustTag(fmt.Sprintf("%%browser-history-%d", int(ct["id"].(float64)))),
		)

	case "tab":
		me.Add(
			ids.MustTag(fmt.Sprintf("%%browser-window_id-%d", int(ct["windowId"].(float64)))),
		)

		me.Add(
			ids.MustTag(fmt.Sprintf("%%browser-tab_id-%d", int(ct["id"].(float64)))),
		)

		v, ok := ct["active"]

		if !ok {
			break
		}

		if b, _ := v.(bool); b {
			me.Add(
				ids.MustTag("%browser-active"),
			)
		}

	case "bookmark":
		me.Add(
			ids.MustTag(fmt.Sprintf("%%browser-bookmark-%d", int(ct["id"].(float64)))),
		)

	}

	return me
}

func (tab browserItemRaw) GetTabId() (id float64, ok bool) {
	switch tab["type"].(string) {
	case "history", "bookmark":
		return
	}

	id, ok = tab["id"].(float64)

	return
}

func (tab browserItemRaw) GetUrl() (u *url.URL, err error) {
	ur := tab["url"]

	if ur == nil {
		err = errors.Errorf("no url: %#v", tab)
		return
	}

	if u, err = url.Parse(ur.(string)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
