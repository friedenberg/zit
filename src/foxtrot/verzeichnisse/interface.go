package verzeichnisse

import (
	"github.com/friedenberg/zit/src/charlie/sha"
)

type IdTransformer func(sha.Sha) string

type Reader interface {
	Begin() (err error)
	ReadRow(string, Row) (err error)
	End() (err error)
}
