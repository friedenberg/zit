package typ

import (
	"github.com/friedenberg/zit/src/delta/kennung"
)

type InlineChecker interface {
	IsInlineTyp(kennung.Typ) bool
}