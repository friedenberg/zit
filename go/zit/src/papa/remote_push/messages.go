package remote_push

import (
	"code.linenisgreat.com/zit-go/src/charlie/gattung"
	"code.linenisgreat.com/zit-go/src/charlie/sha"
)

type messageRequestObjekteData struct {
	Gattung gattung.Gattung
	Sha     sha.Sha
}
