package typ

import (
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
)

func Default() (t typ_akte.V0, k kennung.Typ) {
	k = kennung.MustTyp("md")

	t = typ_akte.Default()

	return
}
