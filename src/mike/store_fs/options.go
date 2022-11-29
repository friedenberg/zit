package store_fs

import (
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/india/zettel"
)

type OptionsReadExternal struct {
	zettel.Format
	Zettelen map[hinweis.Hinweis]zettel.Transacted
}

type CheckoutOptions struct {
	Force bool
	CheckoutMode
	zettel.Format
	Zettelen map[hinweis.Hinweis]zettel.Transacted
}
