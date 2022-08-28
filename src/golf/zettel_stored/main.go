package zettel_stored

import (
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

type Stored struct {
	Sha    sha.Sha
	Zettel zettel.Zettel
}

type Named struct {
	Stored  Stored
	Hinweis hinweis.Hinweis
}

type Transacted struct {
	Named                 Named
	Kopf, Mutter, Schwanz ts.Time
}

type External struct {
	Named
	Path     string
	AktePath string
}
