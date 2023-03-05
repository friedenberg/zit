package objekte

import (
	"github.com/friedenberg/zit/src/delta/kennung"
)

type ObjekteLike interface {
	kennung.Matchable
}
