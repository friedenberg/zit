package zettel_transacted

import (
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/objekte"
	"github.com/friedenberg/zit/src/india/zettel"
)

type zettel_transacted = objekte.Transacted[
	zettel.Zettel,
	*zettel.Zettel,
	hinweis.Hinweis,
	*hinweis.Hinweis,
]

type Transacted = zettel_transacted

// type Zettel = zettel2.Transacted
