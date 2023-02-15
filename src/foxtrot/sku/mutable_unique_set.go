package sku

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type MutableSetUnique = schnittstellen.MutableSet[SkuLike]

func init() {
	gob.Register(
		collections.MakeMutableSet[SkuLike](
			func(s SkuLike) string {
				if s == nil {
					return ""
				}

				return s.GetObjekteSha().String()
			},
		),
	)
}
