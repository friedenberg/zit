package zettel

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type Sku = sku.Transacted[kennung.Hinweis, *kennung.Hinweis]

type Transacted = sku.Transacted[
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
