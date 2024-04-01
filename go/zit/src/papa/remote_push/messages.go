package remote_push

import (
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/delta/sha"
)

type messageRequestObjekteData struct {
	Gattung gattung.Gattung
	Sha     sha.Sha
}
