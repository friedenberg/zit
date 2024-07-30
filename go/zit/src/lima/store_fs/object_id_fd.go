package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type ObjectIdFD struct {
	ids.ObjectId
	fd.FD
}

func (a *ObjectIdFD) String() string {
	return a.FD.String()
}

func (a *ObjectIdFD) Equals(b ObjectIdFD) bool {
	if a.ObjectId.String() != b.ObjectId.String() {
		return false
	}

	if !a.FD.Equals(&b.FD) {
		return false
	}

	return true
}

func (e *ObjectIdFD) GetObjectId() *ids.ObjectId {
	return &e.ObjectId
}
