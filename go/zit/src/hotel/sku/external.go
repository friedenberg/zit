package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type External struct {
	Transacted   Transacted
	ExternalType ids.Type

	// TODO add support for querying the below
	ids.RepoId
	external_state.State
	ExternalObjectId ids.ObjectId
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
	ExternalResetter.ResetWith(b, a)
	return b
}

func (c *External) GetSku() *Transacted {
	return &c.Transacted
}

func (a *External) GetObjectIdLike() ids.IdLike {
	return &a.Transacted.ObjectId
}

func (a *External) GetType() ids.Type {
	return a.Transacted.Metadata.Type
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

func (a *External) GetObjectSha() interfaces.Sha {
	return a.Transacted.Metadata.Sha()
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
