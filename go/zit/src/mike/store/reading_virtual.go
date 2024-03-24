package store

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"syscall"

	"code.linenisgreat.com/chrest"
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/standort"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
)

type chromeTab map[string]interface{}

func (ct chromeTab) Etiketten() kennung.EtikettSet {
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

type VirtualStore struct {
	standort          standort.Standort
	init              sync.Once
	toml_bookmark_typ kennung.Typ
	chromeTabs        map[url.URL]chromeTab
}

func MakeVirtualStore(st standort.Standort) (s *VirtualStore) {
	s = &VirtualStore{
		standort:          st,
		toml_bookmark_typ: kennung.MustTyp("toml-bookmark"),
	}

	return
}

func (vs *VirtualStore) initIfNecessary() (err error) {
	vs.init.Do(
		func() {
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
				if errors.IsErrno(err, syscall.ECONNREFUSED) {
					errors.Err().Print("chrest offline")
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}

			chromeTabsRaw2 := chromeTabsRaw.([]interface{})

			chromeTabs := make(map[url.URL]chromeTab, len(chromeTabsRaw2))

			for _, tabRaw := range chromeTabsRaw2 {
				tab := tabRaw.(map[string]interface{})
				var u *url.URL

				if u, err = url.Parse(tab["url"].(string)); err != nil {
					err = errors.Wrap(err)
					return
				}

				chromeTabs[*u] = tab
			}

			vs.chromeTabs = chromeTabs
		},
	)

	return
}

func (vs *VirtualStore) HydrateOneChrome(sk *sku.Transacted) (err error) {
	// TODO make this more forgiving
	if !sk.Metadatei.Typ.Equals(vs.toml_bookmark_typ) {
		return
	}

	if err = vs.initIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var r sha.ReadCloser

	if r, err = vs.standort.AkteReader(sk.GetAkteSha()); err != nil {
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

	var u *url.URL

	if u, err = url.Parse(tb.Url); err != nil {
		err = errors.Wrap(err)
		return
	}

	tab, ok := vs.chromeTabs[*u]

	if !ok {
		return
	}

	if err = tab.Etiketten().EachPtr(sk.Metadatei.AddEtikettPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
