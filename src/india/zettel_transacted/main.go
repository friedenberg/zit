package zettel_transacted

import (
	"github.com/friedenberg/zit/src/delta/hinweis"
	objekte "github.com/friedenberg/zit/src/delta/objekte"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

type Zettel = objekte.Transacted[
	zettel.Zettel,
	*zettel.Zettel,
	hinweis.Hinweis,
	*hinweis.Hinweis,
]
