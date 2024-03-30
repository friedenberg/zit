package chrome

import (
	"fmt"
	"net/http"
	"net/url"

	"code.linenisgreat.com/chrest"
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/standort"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
)

type Chrome struct {
	typ      kennung.Typ
	standort standort.Standort
	urls     map[url.URL][]item
}

func MakeChrome(s standort.Standort) *Chrome {
	return &Chrome{
		typ:      kennung.MustTyp("toml-bookmark"),
		standort: s,
	}
}

func (c *Chrome) Init() (err error) {
	var chromeTabsRaw interface{}
	var req *http.Request

	if req, err = http.NewRequest("GET", "http://localhost/urls", nil); err != nil {
		err = errors.Wrap(err)
		return
	}

	var chrestConfig chrest.Config

	if err = chrestConfig.Read(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if chromeTabsRaw, err = chrest.AskChrome(chrestConfig, req); err != nil {
		err = errors.Wrap(err)
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

func (c *Chrome) getUrl(sk *sku.Transacted) (u *url.URL, err error) {
	var r sha.ReadCloser

	if r, err = c.standort.AkteReader(sk.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, r)

	var tb sku_fmt.TomlBookmark

	dec := toml.NewDecoder(r)

	if err = dec.Decode(&tb); err != nil {
		err = errors.Wrap(err)
		return
	}

	if u, err = url.Parse(tb.Url); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Chrome) ContainsMatchable(sk *sku.Transacted) bool {
	if !sk.GetTyp().Equals(c.typ) {
		return false
	}

	u, err := c.getUrl(sk)
	errors.PanicIfError(err)

	ts, ok := c.urls[*u]

	if !ok {
		return false
	}

	for _, t := range ts {
		es := t.Etiketten()

		if err = t.Etiketten().EachPtr(sk.Metadatei.AddEtikettPtr); err != nil {
			errors.PanicIfError(err)
		}

		ex := kennung.ExpandMany(es, expansion.ExpanderRight)

		if err = ex.EachPtr(sk.Metadatei.Verzeichnisse.AddEtikettExpandedPtr); err != nil {
			errors.PanicIfError(err)
		}
	}

	return true
}

type item map[string]interface{}

func (ct item) Etiketten() kennung.EtikettSet {
	me := kennung.MakeEtikettMutableSet()

	switch ct["type"].(string) {
	case "history":
		me.Add(
			kennung.MustEtikett(fmt.Sprintf("%%chrome-history-%d", int(ct["id"].(float64)))),
		)

	case "tab":
		me.Add(
			kennung.MustEtikett(fmt.Sprintf("%%chrome-window_id-%d", int(ct["windowId"].(float64)))),
		)

		me.Add(
			kennung.MustEtikett(fmt.Sprintf("%%chrome-tab_id-%d", int(ct["id"].(float64)))),
		)

		v, ok := ct["active"]

		if !ok {
			break
		}

		if b, _ := v.(bool); b {
			me.Add(
				kennung.MustEtikett("%chrome-active"),
			)
		}

	case "bookmark":
		me.Add(
			kennung.MustEtikett(fmt.Sprintf("%%chrome-bookmark-%d", int(ct["id"].(float64)))),
		)

	}

	return me
}
