package chrome

import (
	"code.linenisgreat.com/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

func etiketten(sk *sku.Transacted) kennung.EtikettSet {
	return kennung.ExpandMany(sk.Metadatei.GetEtiketten(), expansion.ExpanderRight)
}
