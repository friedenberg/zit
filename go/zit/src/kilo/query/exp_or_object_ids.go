package query

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

type objectIds struct {
	internal map[string]ObjectId
	external map[string]sku.ExternalObjectId
}

func (oids *objectIds) IsEmpty() bool {
	if len(oids.internal) > 0 {
		return false
	}

	if len(oids.external) > 0 {
		return false
	}

	return true
}

type expOrObjectIds struct {
	isObjectIds bool
	Exp         Exp
	objectIds
}

func (e *expOrObjectIds) IsEmpty() bool {
	if len(e.Exp.Children) > 0 {
		return false
	}

	if !e.objectIds.IsEmpty() {
		return false
	}

	return true
}

func (e *expOrObjectIds) reduce(b *buildState) (err error) {
	if err = e.Exp.reduce(b); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
