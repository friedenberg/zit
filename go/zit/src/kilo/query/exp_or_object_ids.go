package query

import "code.linenisgreat.com/zit/go/zit/src/juliett/sku"

type expOrObjectIds struct {
	isObjectIds bool
	Exp
	objectIds
}

type objectIds struct {
	internal map[string]ObjectId
	external map[string]sku.ExternalObjectId
}
