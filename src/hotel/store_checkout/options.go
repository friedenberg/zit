package store_checkout

import (
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
)

type OptionsReadExternal struct {
	zettel.Format
	Zettelen map[hinweis.Hinweis]stored_zettel.Transacted
}

type CheckoutOptions struct {
	CheckoutMode
	zettel.Format
	Zettelen map[hinweis.Hinweis]stored_zettel.Transacted
}
