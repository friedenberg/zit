package remote_pull

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
)

type messageServeData struct {
	Size int64
}

type messageServeAkte messageServeData
type messageServeObjekte struct {
	Gattung gattung.Gattung
	messageServeData
}
