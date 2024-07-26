package store

import (
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) addToTomlIndexIfNecessary(
	t *sku.Transacted,
	o sku.CommitOptions,
) (err error) {
	if !t.Metadata.Type.IsToml() {
		return
	}

	return
}
