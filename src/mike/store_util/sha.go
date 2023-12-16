package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func (s *common) CalculateAndSetShaTransacted(sk *sku.Transacted) (err error) {
	return s.CalculateAndSetShaSkuLike(sk)
}

func (s *common) CalculateAndSetShaSkuLike(sk sku.SkuLike) (err error) {
	if err = sku.CalculateAndSetSha(
		sk,
		s.persistentMetadateiFormat,
		s.options,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
