package store_fs

import (
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/hotel/zettel"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
)

type OptionsReadExternal struct {
	zettel.Format
	Zettelen map[hinweis.Hinweis]zettel_transacted.Zettel
}

type CheckoutOptions struct {
	Force bool
	CheckoutMode
	zettel.Format
	Zettelen map[hinweis.Hinweis]zettel_transacted.Zettel
}
