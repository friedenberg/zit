package typ

import (
	"github.com/friedenberg/zit/src/charlie/kennung"
	"github.com/friedenberg/zit/src/delta/konfig"
)

type Kennung = kennung.Typ

func IsInlineAkte(t Kennung, k konfig.Konfig) (isInline bool) {
	ts := t.String()
	tc := k.GetTyp(ts)

	if tc == nil {
		return
	}

	isInline = tc.InlineAkte

	return
}
