package remote_push

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type messageRequestObjekteData struct {
	Gattung genres.Genre
	Sha     sha.Sha
}
