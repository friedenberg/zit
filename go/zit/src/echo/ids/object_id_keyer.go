package ids

import "code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"

type ObjectIdGetter interface {
	GetObjectId() *ObjectId
}

type ObjectIdKeyer[
	T any,
	TPtr interface {
		interfaces.Ptr[T]
		ObjectIdGetter
	},
] struct{}

func (sk ObjectIdKeyer[T, TPtr]) GetKey(e TPtr) string {
	return e.GetObjectId().String()
}
