package typ

import (
	"github.com/friedenberg/zit/src/foxtrot/kennung"
)

type InlineChecker interface {
	IsInlineTyp(kennung.Typ) bool
}
