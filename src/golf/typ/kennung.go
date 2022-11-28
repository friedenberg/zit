package typ

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/konfig"
)

func IsInlineAkte(t kennung.Typ, k konfig.Konfig) (isInline bool) {
	ts := t.String()
	tc := k.Transacted.Objekte.GetTyp(ts)

	if tc == nil {
		return
	}

	isInline = tc.Typ.Akte.InlineAkte

	return
}
