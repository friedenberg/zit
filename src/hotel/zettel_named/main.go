package zettel_named

import (
	"github.com/friedenberg/zit/src/delta/hinweis"
	objekte "github.com/friedenberg/zit/src/delta/objekte"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

type Zettel = objekte.Named[
	zettel.Zettel,
	*zettel.Zettel,
	hinweis.Hinweis,
	*hinweis.Hinweis,
]
