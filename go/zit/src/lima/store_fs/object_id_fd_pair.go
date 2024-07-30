package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type ObjectIdFDPair struct {
	ObjectId ids.ObjectId
	FDs      FDPair
}

func (a *ObjectIdFDPair) String() string {
	return a.ObjectId.String()
}

func (a *ObjectIdFDPair) Equals(b ObjectIdFDPair) bool {
	if a.ObjectId.String() != b.ObjectId.String() {
		return false
	}

	if !a.FDs.Equals(&b.FDs) {
		return false
	}

	return true
}

func (e *ObjectIdFDPair) GetObjectId() *ids.ObjectId {
	return &e.ObjectId
}

func (e *ObjectIdFDPair) GetFDs() *FDPair {
	return &e.FDs
}

func (e *ObjectIdFDPair) GetObjectFD() *fd.FD {
	return &e.FDs.Object
}

func (e *ObjectIdFDPair) GetBlobFD() *fd.FD {
	return &e.FDs.Blob
}
