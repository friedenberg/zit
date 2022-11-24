package zettel_transacted

import (
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/objekte"
	"github.com/friedenberg/zit/src/hotel/zettel"
)

type Zettel = objekte.Transacted[
	zettel.Zettel,
	*zettel.Zettel,
	hinweis.Hinweis,
	*hinweis.Hinweis,
]
