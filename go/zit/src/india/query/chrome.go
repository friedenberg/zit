package query

import (
	"fmt"
	"net/http"
	"net/url"

	"code.linenisgreat.com/chrest"
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/standort"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
)

type Chrome struct {
	typ      kennung.Typ
	standort standort.Standort
	tabs     map[url.URL]tab
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

	if req, err = http.NewRequest("GET", "http://localhost/tabs", nil); err != nil {
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

	chromeTabs := make(map[url.URL]tab, len(chromeTabsRaw2))

	for _, tabRaw := range chromeTabsRaw2 {
		tab := tabRaw.(map[string]interface{})
		var u *url.URL

		if u, err = url.Parse(tab["url"].(string)); err != nil {
			err = errors.Wrap(err)
			return
		}

		chromeTabs[*u] = tab
	}

	c.tabs = chromeTabs

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

func (c *Chrome) String() string {
	return "%chrome"
}

func (c *Chrome) Each(_ schnittstellen.FuncIter[Matcher]) error {
	return nil
}

func (c *Chrome) ContainsMatchable(sk *sku.Transacted) bool {
	if !sk.GetTyp().Equals(c.typ) {
		return false
	}

	u, err := c.getUrl(sk)
	errors.PanicIfError(err)

	t, ok := c.tabs[*u]

	if !ok {
		return false
	}

	if err = t.Etiketten().EachPtr(sk.Metadatei.AddEtikettPtr); err != nil {
		errors.PanicIfError(err)
	}

	return true
}

type tab map[string]interface{}

func (ct tab) Etiketten() kennung.EtikettSet {
	me := kennung.MakeEtikettMutableSet()

	me.Add(
		kennung.MustEtikett(fmt.Sprintf("%%chrome-window_id-%d", int(ct["windowId"].(float64)))),
	)

	me.Add(
		kennung.MustEtikett(fmt.Sprintf("%%chrome-tab_id-%d", int(ct["id"].(float64)))),
	)

	if ct["active"].(bool) {
		me.Add(
			kennung.MustEtikett("%chrome-active"),
		)
	}

	return me
}
