package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

type External struct {
	Transacted
  Fields map[string]Field
}

func (t *External) GetSkuExternalLike() ExternalLike {
	return t
}

func (a *External) Clone() ExternalLike {
	b := GetExternalPool().Get()
	TransactedResetter.ResetWith(&b.Transacted, &a.Transacted)
	return b
}

func (c *External) GetSku() *Transacted {
	return &c.Transacted
}

func (a *External) GetObjectIdLike() ids.IdLike {
	return &a.ObjectId
}

func (a *External) GetMetadatei() *object_metadata.Metadata {
	return &a.Metadata
}

func (a *External) GetGenre() interfaces.Genre {
	return a.ObjectId.GetGenre()
}

func (a *External) String() string {
	return fmt.Sprintf(
		". %s %s %s %s",
		a.GetGenre(),
		a.GetObjectIdLike(),
		a.GetObjectSha(),
		a.GetBlobSha(),
	)
}

func (a *External) GetBlobSha() interfaces.Sha {
	return &a.Metadata.Blob
}

func (a *External) SetBlobSha(v interfaces.Sha) (err error) {
	if err = a.Metadata.Blob.SetShaLike(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *External) AsTransacted() (b Transacted) {
	b = a.Transacted

	return
}

func (o *External) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGenre(), o.GetObjectIdLike())
}
