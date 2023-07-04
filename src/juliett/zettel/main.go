package zettel

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/golf/sku"
)

type Sku = sku.Transacted[kennung.Hinweis, *kennung.Hinweis]

type Transacted = objekte.Transacted[
	Objekte,
	*Objekte,
	kennung.Hinweis,
	*kennung.Hinweis,
]

type External = objekte.External[
	Objekte,
	*Objekte,
	kennung.Hinweis,
	*kennung.Hinweis,
]

type CheckedOut = objekte.CheckedOut[
	Objekte,
	*Objekte,
	kennung.Hinweis,
	*kennung.Hinweis,
]
