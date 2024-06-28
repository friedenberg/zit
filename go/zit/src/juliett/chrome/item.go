package chrome

import (
	"fmt"
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/bezeichnung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type item map[string]interface{}

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

func (tab item) GetTai() (t kennung.Tai, err error) {
	date, ok := tab["date"].(string)

	if !ok {
		err = errors.Errorf("expected string but got %T, %q", tab["date"], tab["date"])
		return
	}

	if err = t.SetFromRFC3339(date); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (tab item) GetBezeichnung() (b bezeichnung.Bezeichnung, err error) {
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

func (tab item) GetTyp() (t kennung.Typ, err error) {
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
