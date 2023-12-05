package store_verzeichnisse

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type KennungShaMap map[string]*sha.Sha

func (ksm KennungShaMap) ModifyMutter(z *sku.Transacted) (err error) {
	k := z.GetKennung()
	old, ok := ksm[k.String()]

	if !ok {
		return
	}

	if err = z.GetMetadatei().Verzeichnisse.Mutter.SetShaLike(old); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ksm KennungShaMap) SaveSha(z *sku.Transacted) (err error) {
	k := z.GetKennung()
	var sh sha.Sha

	if err = sh.SetShaLike(&z.GetMetadatei().Verzeichnisse.Sha); err != nil {
		err = errors.Wrap(err)
		return
	}

	ksm[k.String()] = &sh

	return
}
