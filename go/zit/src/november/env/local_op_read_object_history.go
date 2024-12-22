package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (env *Local) ReadObjectHistory(
	oid *ids.ObjectId,
) (skus []*sku.Transacted, err error) {
	streamIndex := env.GetStore().GetStreamIndex()

	if skus, err = streamIndex.ReadManyObjectId(
		oid.String(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
