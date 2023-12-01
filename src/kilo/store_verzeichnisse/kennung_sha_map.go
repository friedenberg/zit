package store_verzeichnisse

import (
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type KennungShaMap map[string]sha.Sha

func (ksm KennungShaMap) ModifyMutter(z *sku.Transacted) (err error) {
	k := z.GetKennung()
	old, ok := ksm[k.String()]

	if !ok {
		return
	}

	z.GetMetadatei().Verzeichnisse.Mutter = old

	return
}

func (ksm KennungShaMap) SaveSha(z *sku.Transacted) (err error) {
	k := z.GetKennung()
	ksm[k.String()] = z.GetMetadatei().Verzeichnisse.Sha

	return
}
