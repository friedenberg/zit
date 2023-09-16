package store_verzeichnisse

import (
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/india/transacted"
)

type KennungShaMap map[kennung.Hinweis]sha.Sha

func (ksm KennungShaMap) ModifyMutter(z *transacted.Zettel) (err error) {
	k := z.GetKennung()
	old, ok := ksm[k]

	if !ok {
		return
	}

	z.GetMetadateiPtr().Verzeichnisse.Mutter[0] = old

	return
}

func (ksm KennungShaMap) SaveSha(z *transacted.Zettel) (err error) {
	k := z.GetKennung()
	ksm[k] = z.GetMetadatei().Verzeichnisse.Sha

	return
}
