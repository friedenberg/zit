package stored_zettel

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
	Stored
	Hinweis hinweis.Hinweis
}

type Transacted struct {
	Named
	Head, Tail ts.Time
}

type External struct {
	Path     string
	AktePath string
	Hinweis  hinweis.Hinweis
	Sha      sha.Sha
	Zettel   zettel.Zettel
}

type CheckedOut struct {
	Internal Transacted
	External External
}
