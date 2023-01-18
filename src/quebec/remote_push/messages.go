package remote_push

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/foxtrot/sha"
)

type messageRequestObjekteData struct {
	Gattung gattung.Gattung
	Sha     sha.Sha
}
