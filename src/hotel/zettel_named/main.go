package zettel_named

import (
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/objekte_stored"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

type Zettel = objekte_stored.Named[
	zettel.Zettel,
	*zettel.Zettel,
	hinweis.Hinweis,
	*hinweis.Hinweis,
]
