package chrome

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
)

// TODO make more specific
type item map[string]interface{}

func (item item) WriteToMetadatei(m *metadatei.Metadatei) (err error) {
	if m.Tai, err = item.GetTai(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if m.Typ, err = item.GetTyp(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if m.Bezeichnung, err = item.GetBezeichnung(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var u *url.URL

	if u, err = item.GetUrl(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var e ids.Tag

	els := strings.Split(u.Host, ".")
	slices.Reverse(els)
	host := strings.Join(els, "-")

	if err = e.Set("zz-site-" + host); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = m.AddEtikettPtr(&e); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO add etiketten

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

func (tab item) GetBezeichnung() (b descriptions.Description, err error) {
	t, ok := tab["title"].(string)

	if !ok {
		err = errors.Errorf("expected string but got %T, %q", tab["title"], tab["title"])
		return
	}

	if err = b.Set(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (tab item) GetTyp() (t ids.Type, err error) {
	ty, ok := tab["type"].(string)

	if !ok {
		err = errors.Errorf("expected string but got %T, %q", tab["type"], tab["type"])
		return
	}

	if err = t.Set("chrome-" + ty); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ct item) Etiketten() ids.TagSet {
	me := ids.MakeTagMutableSet()

	switch ct["type"].(string) {
	case "history":
		me.Add(
			ids.MustTag(fmt.Sprintf("%%chrome-history-%d", int(ct["id"].(float64)))),
		)

	case "tab":
		me.Add(
			ids.MustTag(fmt.Sprintf("%%chrome-window_id-%d", int(ct["windowId"].(float64)))),
		)

		me.Add(
			ids.MustTag(fmt.Sprintf("%%chrome-tab_id-%d", int(ct["id"].(float64)))),
		)

		v, ok := ct["active"]

		if !ok {
			break
		}

		if b, _ := v.(bool); b {
			me.Add(
				ids.MustTag("%chrome-active"),
			)
		}

	case "bookmark":
		me.Add(
			ids.MustTag(fmt.Sprintf("%%chrome-bookmark-%d", int(ct["id"].(float64)))),
		)

	}

	return me
}
