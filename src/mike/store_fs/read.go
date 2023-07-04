package store_fs

import (
	"fmt"
	"path"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/golf/objekte"
)

func (s *Store) FileExtensionForGattung(
	gg schnittstellen.GattungGetter,
) string {
	return s.erworben.FileExtensions.GetFileExtensionForGattung(gg)
}

func (s *Store) PathForTransactedLike(tl objekte.TransactedLike) string {
	return path.Join(
		s.Cwd(),
		fmt.Sprintf(
			"%s.%s",
			tl.GetSkuLike().GetId(),
			s.FileExtensionForGattung(tl.GetSkuLike()),
		),
	)
}
