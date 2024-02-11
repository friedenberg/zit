package typ

import (
	"code.linenisgreat.com/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

func Default() (t typ_akte.V0, k kennung.Typ) {
	k = kennung.MustTyp("md")

	t = typ_akte.Default()

	return
}
