package ids

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type IndexedLike struct {
	Int int
	ObjectId
	SchwanzenCount int
	Count          int
}

func MakeIndexed(k interfaces.ObjectId) (i *IndexedLike) {
	i = &IndexedLike{}
	i.ResetWithObjectId(k)
	return
}

func (i *IndexedLike) ResetWithObjectId(k interfaces.ObjectId) {
	errors.PanicIfError(i.ObjectId.SetWithIdLike(k))
}

func (z *IndexedLike) GetInt() int {
	return 0
}

func (z *IndexedLike) GetObjectId() interfaces.ObjectId {
	return &z.ObjectId
}

func (k *IndexedLike) GetSchwanzenCount() int {
	return k.SchwanzenCount
}

func (k *IndexedLike) GetCount() int {
	return k.Count
}

func (z *IndexedLike) Reset() {
	z.ObjectId.Reset()
	z.SchwanzenCount = 0
	z.Count = 0
}
