package sku

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/delta/collections"
)

type MutableSetUnique = collections.MutableSet[*Sku]

func init() {
	gob.Register(
		collections.MakeMutableSet[*Sku](
			func(s *Sku) string {
				if s == nil {
					return ""
				}

				return s.Sha.String()
			},
		),
	)
}
