package store_fs

import (
	"fmt"
	"path"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func (s *Store) FileExtensionForGattung(
	gg schnittstellen.GattungGetter,
) string {
	return s.GetKonfig().FileExtensions.GetFileExtensionForGattung(gg)
}

func (s *Store) PathForTransactedLike(tl sku.SkuLike) string {
	return path.Join(
		s.Cwd(),
		fmt.Sprintf(
			"%s.%s",
			tl.GetSkuLike().GetKennungLike(),
			s.FileExtensionForGattung(tl.GetSkuLike()),
		),
	)
}
