package objekte_index

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

type IndexGenericPresence interface {
	ExistsAkteSha(schnittstellen.Sha) bool
	ExistsObjekteSha(schnittstellen.Sha) bool
}
