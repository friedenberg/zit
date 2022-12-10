package typ

import (
	"github.com/friedenberg/zit/src/foxtrot/kennung"
)

type typGetter interface {
	GetTyp(kennung.Typ) *Transacted
}

func IsInlineAkte(t kennung.Typ, k typGetter) (isInline bool) {
	tc := k.GetTyp(t)

	if tc == nil {
		return
	}

	isInline = tc.Objekte.Akte.InlineAkte

	return
}
