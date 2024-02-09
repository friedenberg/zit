package store_objekten

import (
	"github.com/friedenberg/zit/src/bravo/objekte_mode"
	"github.com/friedenberg/zit/src/hotel/sku"
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
