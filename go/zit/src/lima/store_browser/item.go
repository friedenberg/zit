package store_browser

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

// TODO make more specific
type item map[string]interface{}

func (item item) WriteToMetadata(m *object_metadata.Metadata) (err error) {
	if m.Tai, err = item.GetTai(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if m.Type, err = item.GetType(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if m.Description, err = item.GetDescription(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var u *url.URL

	if u, err = item.GetUrl(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var e ids.Tag

	els := strings.Split(u.Hostname(), ".")
	slices.Reverse(els)

	if els[0] == "www" {
		els = els[1:]
	}

	host := strings.Join(els, "-")

	if len(host) > 0 {
		if err = e.Set("zz-site-" + host); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = m.AddTagPtr(&e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (tab item) WriteToObjectId(oi *ids.ObjectId) (err error) {
	ty, ok := tab["type"].(string)

	if !ok {
		err = errors.Errorf("expected string but got %T, %q", tab["type"], tab["type"])
		return
	}

	var id string
	id, ok = tab["id"].(string)

	if !ok {
		err = errors.Errorf("unsupported id format: %#v", id)
		return
	}

	if err = oi.SetRaw(fmt.Sprintf("%s-%s", ty, id)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (tab item) GetTabId() (id float64, ok bool) {
	switch tab["type"].(string) {
	case "history", "bookmark":
		return
	}

	id, ok = tab["id"].(float64)

	return
}

func (tab item) GetUrl() (u *url.URL, err error) {
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

func (tab item) GetTai() (t ids.Tai, err error) {
	switch date := tab["date"].(type) {
	case nil:
		t = ids.NowTai()

	case string:
		if err = t.SetFromRFC3339(date); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = errors.Errorf("expected string but got %T, %q", tab["date"], tab["date"])
		return
	}

	return
}

func (tab item) GetDescription() (b descriptions.Description, err error) {
	switch t := tab["title"].(type) {
	case nil:
	case string:
		if err = b.Set(t); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = errors.Errorf("expected string but got %T, %q", t, t)
		return
	}

	return
}

func (tab item) GetType() (t ids.Type, err error) {
	ty, ok := tab["type"].(string)

	if !ok {
		err = errors.Errorf("expected string but got %T, %q", tab["type"], tab["type"])
		return
	}

	if err = t.Set("browser-" + ty); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ct item) GetTagSet() ids.TagSet {
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
