package store_fs

import (
	"github.com/friedenberg/zit/src/echo/hinweis"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type OptionsReadExternal struct {
	Zettelen map[hinweis.Hinweis]zettel.Transacted
}

type CheckoutOptions struct {
	Force bool
	CheckoutMode
	Zettelen map[hinweis.Hinweis]zettel.Transacted
}
