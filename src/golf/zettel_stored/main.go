package zettel_stored

import (
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

type Stored struct {
	Sha    sha.Sha
	Zettel zettel.Zettel
}

func (zs *Stored) Reset() {
	zs.Sha = sha.Sha{}
	zs.Zettel.Reset()
}
