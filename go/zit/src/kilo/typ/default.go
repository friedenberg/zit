package typ

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

func Default() (t type_blobs.V0, k ids.Type) {
	k = ids.MustType("md")

	t = type_blobs.Default()

	return
}
