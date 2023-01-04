package sku

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/delta/collections"
)

type MutableSetUnique = collections.MutableSet[SkuLike]

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
