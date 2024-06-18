package chrome

import (
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO replace with etiketten_path.Etiketten
func etiketten(sk *sku.Transacted) kennung.EtikettSet {
	return kennung.ExpandMany(sk.Metadatei.GetEtiketten(), expansion.ExpanderRight)
}
