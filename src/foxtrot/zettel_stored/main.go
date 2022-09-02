package zettel_stored

import (
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/zettel"
)

type Stored struct {
	Sha    sha.Sha
	Zettel zettel.Zettel
}
