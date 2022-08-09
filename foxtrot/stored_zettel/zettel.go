package stored_zettel

import (
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/ts"
)

type Stored struct {
	Mutter sha.Sha
	Kinder sha.Sha
	Sha    sha.Sha
	Zettel zettel.Zettel
}

type Named struct {
	Stored
	Hinweis hinweis.Hinweis
}

type Transacted struct {
	Named
	ts.Time
}

type External struct {
	Path     string
	AktePath string
	Hinweis  hinweis.Hinweis
	Sha      sha.Sha
	Zettel   zettel.Zettel
}

type CheckedOut struct {
	//TODO change to Transacted
	Internal Named
	External External
}
