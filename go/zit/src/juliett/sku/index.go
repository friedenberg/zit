package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type (
	IndexObject interface {
		ReadOneObjectIdTai(
			k interfaces.ObjectId,
			t ids.Tai,
		) (sk *Transacted, err error)

		ReadManyObjectId(
			id interfaces.ObjectId,
		) (skus []*Transacted, err error)

		ReadOneObjectId(
			oid interfaces.ObjectId,
			sk *Transacted,
		) (err error)

		ObjectExists(
			id *ids.ObjectId,
		) (err error)

		ReadManySha(
			sh *sha.Sha,
		) (skus []*Transacted, err error)

		ReadOneSha(
			sh *sha.Sha,
			sk *Transacted,
		) (err error)
	}
)
