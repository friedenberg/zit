package store_fs

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type OptionsReadExternal struct {
	Zettelen map[kennung.Hinweis]zettel.Transacted
}

type CheckoutOptions struct {
	Force        bool
	CheckoutMode objekte.CheckoutMode
	Zettelen     map[kennung.Hinweis]zettel.Transacted
}
