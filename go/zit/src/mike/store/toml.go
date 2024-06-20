package store

import (
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) addToTomlIndexIfNecessary(
	t *sku.Transacted,
	o ObjekteOptions,
) (err error) {
	if !t.Metadatei.Typ.IsToml() {
		return
	}

	return
}
