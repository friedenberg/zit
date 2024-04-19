package chrome

import (
	"fmt"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type item map[string]interface{}

func (ct item) HydrateSku(sk *sku.Transacted) (err error) {
	if date, ok := ct["date"].(string); ok {
		if err = sk.Metadatei.Tai.SetFromRFC3339(date); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = sk.Metadatei.Bezeichnung.Set(ct["title"].(string)); err != nil {
		err = errors.Wrap(err)
		return
	}

	t := ct["type"].(string)
	sk.Metadatei.Typ = kennung.MustTyp("chrome-" + t)

	if err = sk.Kennung.Set(
		fmt.Sprintf("%%%d", int(ct["id"].(float64))),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	me := sk.GetMetadatei().GetEtikettenMutable()

	switch t {
	case "history":

	case "tab":
		me.Add(
			kennung.MustEtikett(fmt.Sprintf("%%chrome-window_id-%d", int(ct["windowId"].(float64)))),
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
