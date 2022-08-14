package stored_zettel

import (
	"github.com/friedenberg/zit/charlie/sha"
	"github.com/friedenberg/zit/delta/hinweis"
	"github.com/friedenberg/zit/delta/ts"
	"github.com/friedenberg/zit/foxtrot/zettel"
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
	//TODO change to Transacted
	Internal Transacted
	External External
}
