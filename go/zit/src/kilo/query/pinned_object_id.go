package query

import "code.linenisgreat.com/zit/go/zit/src/echo/ids"

type pinnedObjectId struct {
	ids.Sigil
	ObjectId
}
