package store_objekten

import "github.com/friedenberg/zit/src/foxtrot/kennung"

func (s Store) Etiketten() (es []kennung.Etikett, err error) {
	return s.zettelStore.indexEtiketten.allEtiketten()
}