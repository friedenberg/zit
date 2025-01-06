package ids

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
)

func MakeExternalObjectId(g genres.Genre, value string) *ExternalObjectId {
	return &ExternalObjectId{
		value: value,
		genre: g,
	}
}

type ExternalObjectId struct {
	value string
	genre genres.Genre
}

func (eoid *ExternalObjectId) GetExternalObjectId() ExternalObjectIdLike {
	return eoid
}

func (eoid *ExternalObjectId) GetGenre() interfaces.Genre {
	return eoid.genre
}

func (eoid *ExternalObjectId) IsEmpty() bool {
	return eoid.value == ""
}

func (eoid *ExternalObjectId) String() string {
	return eoid.value
}
