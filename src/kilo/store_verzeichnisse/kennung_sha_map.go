package store_verzeichnisse

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type KennungShaMap map[string]*sha.Sha

func (ksm KennungShaMap) ReadMutter(z *sku.Transacted) (err error) {
	k := z.GetKennung()
	old := ksm[k.String()]

	if old.IsNull() {
		return
	}

	if err = z.GetMetadatei().Mutter.SetShaLike(old); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ksm KennungShaMap) SaveSha(z *sku.Transacted) (err error) {
	k := z.GetKennung()
	var sh sha.Sha

	if err = sh.SetShaLike(&z.GetMetadatei().Sha); err != nil {
		err = errors.Wrap(err)
		return
	}

	ksm[k.String()] = &sh

	return
}
