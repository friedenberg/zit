package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

type External struct {
	Transacted TransactedWithFields

	// TODO add support for querying the below
	ids.RepoId
	external_state.State
	ExternalObjectId ids.ObjectId
	ExternalType     ids.Type
}

func (t *External) GetRepoId() ids.RepoId {
	return t.RepoId
}

func (t *External) GetObjectId() *ids.ObjectId {
	return &t.Transacted.ObjectId
}

func (t *External) GetExternalObjectId() ids.ExternalObjectId {
	return &t.ExternalObjectId
}

func (t *External) GetSkuExternalLike() ExternalLike {
	return &t.Transacted
}

func (t *External) GetExternalState() external_state.State {
	return external_state.Unknown
}

func (a *External) Clone() ExternalLike {
	b := GetExternalPool().Get()
	TransactedResetter.ResetWith(b.GetSku(), a.GetSku())
	return b
}

func (c *External) GetSku() *Transacted {
	return &c.Transacted.Transacted
}

func (a *External) GetObjectIdLike() ids.IdLike {
	return &a.Transacted.ObjectId
}

func (a *External) GetMetadatei() *object_metadata.Metadata {
	return &a.Transacted.Metadata
}

func (a *External) GetGenre() interfaces.Genre {
	return a.Transacted.ObjectId.GetGenre()
}

func (a *External) String() string {
	return fmt.Sprintf(
		". %s %s %s %s",
		a.GetGenre(),
		a.GetObjectIdLike(),
		a.Transacted.GetObjectSha(),
		a.GetBlobSha(),
	)
}

func (a *External) GetBlobSha() interfaces.Sha {
	return &a.Transacted.Metadata.Blob
}

func (a *External) SetBlobSha(v interfaces.Sha) (err error) {
	if err = a.Transacted.Metadata.Blob.SetShaLike(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (o *External) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGenre(), o.GetObjectIdLike())
}
