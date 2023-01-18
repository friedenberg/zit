package typ

import (
	"github.com/friedenberg/zit/src/echo/kennung"
)

type InlineChecker interface {
	IsInlineTyp(kennung.Typ) bool
}
