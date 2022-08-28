package store_working_directory

import (
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	zettel_stored "github.com/friedenberg/zit/src/golf/zettel_stored"
)

type OptionsReadExternal struct {
	zettel.Format
	Zettelen map[hinweis.Hinweis]zettel_stored.Transacted
}

type CheckoutOptions struct {
	CheckoutMode
	zettel.Format
	Zettelen map[hinweis.Hinweis]zettel_stored.Transacted
}
