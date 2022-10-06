package zettel_stored

import (
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/zettel"
)

type Stored struct {
	Sha    sha.Sha
	Zettel zettel.Zettel
}

func (zs *Stored) Reset() {
  zs.Sha = sha.Sha{}
  zs.Zettel.Reset()
}
