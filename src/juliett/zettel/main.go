package zettel

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type CheckedOut = objekte.CheckedOut[
	kennung.Hinweis,
	*kennung.Hinweis,
]
