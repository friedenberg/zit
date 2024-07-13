package typ

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/type_blob"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

func Default() (t type_blob.V0, k ids.Type) {
	k = ids.MustType("md")

	t = type_blob.Default()

	return
}
