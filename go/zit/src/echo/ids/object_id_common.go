package ids

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
)

func GetObjectIdPool() interfaces.Pool[ObjectId, *ObjectId] {
	return getObjectIdPool2()
}

type ObjectId = objectId2

type IdParts struct {
	Middle              byte
	RepoId, Left, Right *catgut.String
}

var ErrFDNotId = errors.New("not a id file")
