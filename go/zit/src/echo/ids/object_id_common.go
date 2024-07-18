package ids

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
)


func GetObjectIdPool() interfaces.Pool[ObjectId, *ObjectId] {
	return getObjectIdPool()
}

type ObjectId = objectId

func MakeId(v string) (IdLikePtr, error) {
	k := &ObjectId{
		g: genres.Unknown,
	}

	return k, k.Set(v)
}

type IdParts struct {
	Middle              byte
	RepoId, Left, Right *catgut.String
}

var ErrFDNotId = errors.New("not a id file")
