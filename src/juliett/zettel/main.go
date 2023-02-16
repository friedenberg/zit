package zettel

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type Stored = objekte.Stored[
	Objekte,
	*Objekte,
	kennung.Hinweis,
	*kennung.Hinweis,
]

type Sku = sku.Transacted[kennung.Hinweis, *kennung.Hinweis]

type Transacted = objekte.Transacted[
	Objekte,
	*Objekte,
	kennung.Hinweis,
	*kennung.Hinweis,
	Verzeichnisse,
	*Verzeichnisse,
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
	Verzeichnisse,
	*Verzeichnisse,
]
