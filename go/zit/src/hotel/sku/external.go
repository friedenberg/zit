package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type External struct {
	Transacted Transacted
	ExternalInfo
}

func (t *External) GetObjectId() *ids.ObjectId {
	return &t.Transacted.ObjectId
}

func (t *External) GetSkuExternalLike() ExternalLike {
	return &t.Transacted
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
		a.Transacted.GetBlobSha(),
	)
}

func (o *External) GetKey() string {
	return fmt.Sprintf("%s.%s", o.GetGenre(), o.GetObjectIdLike())
}
