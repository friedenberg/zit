package store_verzeichnisse

import (
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/india/transacted"
)

type KennungShaMap map[string]sha.Sha

func (ksm KennungShaMap) ModifyMutter(z *transacted.Zettel) (err error) {
	k := z.GetKennung()
	old, ok := ksm[k.String()]

	if !ok {
		return
	}

	z.GetMetadateiPtr().Verzeichnisse.Mutter = old

	return
}

func (ksm KennungShaMap) SaveSha(z *transacted.Zettel) (err error) {
	k := z.GetKennung()
	ksm[k.String()] = z.GetMetadatei().Verzeichnisse.Sha

	return
}
