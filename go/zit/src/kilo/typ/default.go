package typ

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/type_blob"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

func Default() (t type_blob.V0, k kennung.Type) {
	k = kennung.MustType("md")

	t = type_blob.Default()

	return
}
