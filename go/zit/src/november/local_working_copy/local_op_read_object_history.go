package local_working_copy

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func (env *Repo) ReadObjectHistory(
	oid *ids.ObjectId,
) (skus []*sku.Transacted, err error) {
	streamIndex := env.GetStore().GetStreamIndex()

	if skus, err = streamIndex.ReadManyObjectId(
		oid,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
