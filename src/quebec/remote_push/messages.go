package remote_push

import (
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
)

type messageRequestObjekteData struct {
	Gattung gattung.Gattung
	Sha     sha.Sha
}
