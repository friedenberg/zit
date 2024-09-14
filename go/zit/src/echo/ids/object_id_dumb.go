package ids

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
)

type DumbObjectId struct {
	Value string
	Genre genres.Genre
}

func (a *DumbObjectId) IsEmpty() bool {
	return a.Value == ""
}

func (a *DumbObjectId) GetExternalObjectId() ExternalObjectId {
	return a
}

func (a *DumbObjectId) CloneExternalObjectId() ExternalObjectId {
	b := *a
	return &b
}

func (a *DumbObjectId) String() string {
	return a.Value
}

func (a *DumbObjectId) GetGenre() interfaces.Genre {
	return a.Genre
}
