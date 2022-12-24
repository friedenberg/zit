package zettel

import (
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type Stored = objekte.Stored[
	Objekte,
	*Objekte,
	hinweis.Hinweis,
	*hinweis.Hinweis,
]

type Sku = sku.Transacted[hinweis.Hinweis, *hinweis.Hinweis]

type Transacted = objekte.Transacted[
	Objekte,
	*Objekte,
	hinweis.Hinweis,
	*hinweis.Hinweis,
	Verzeichnisse,
	*Verzeichnisse,
]
