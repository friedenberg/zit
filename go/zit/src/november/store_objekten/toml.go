package store_objekten

import (
	"code.linenisgreat.com/zit-go/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
)

func (s *Store) addToTomlIndexIfNecessary(
	t *sku.Transacted,
	mode objekte_mode.Mode,
) (err error) {
	if !t.Metadatei.Typ.IsToml() {
		return
	}

	return
}
