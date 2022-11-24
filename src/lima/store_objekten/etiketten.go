package store_objekten

import "github.com/friedenberg/zit/src/delta/kennung"

func (s Store) Etiketten() (es []kennung.Etikett, err error) {
	return s.indexEtiketten.allEtiketten()
}
