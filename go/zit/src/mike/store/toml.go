package store

import (
	"code.linenisgreat.com/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/src/hotel/sku"
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
