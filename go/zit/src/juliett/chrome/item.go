package chrome

import (
	"fmt"

	"code.linenisgreat.com/zit/src/echo/kennung"
)

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
