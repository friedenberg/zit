package typ

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type CheckedOut = objekte.CheckedOut[
	kennung.Typ,
	*kennung.Typ,
]

// TODO-P1 move to konfig
// func GetFileExtension(t *Transacted, agp *Akte) string {
// 	if a.FileExtension != "" {
// 		return a.FileExtension
// 	}

// 	return t.GetKennungLike().String()
// }
