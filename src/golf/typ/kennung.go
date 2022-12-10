package typ

import (
	"github.com/friedenberg/zit/src/echo/kennung"
)

type typGetter interface {
	GetTyp(string) *Transacted
}

func IsInlineAkte(t kennung.Typ, k typGetter) (isInline bool) {
	ts := t.String()
	tc := k.GetTyp(ts)

	if tc == nil {
		return
	}

	isInline = tc.Objekte.Akte.InlineAkte

	return
}
