package store_working_directory

import (
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
)

type OptionsReadExternal struct {
	zettel.Format
	Zettelen map[hinweis.Hinweis]zettel_transacted.Zettel
}

type CheckoutOptions struct {
	CheckoutMode
	zettel.Format
	Zettelen map[hinweis.Hinweis]zettel_transacted.Zettel
}
