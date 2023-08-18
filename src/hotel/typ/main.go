package typ

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type Transacted = objekte.Transacted[
	Akte,
	*Akte,
	kennung.Typ,
	*kennung.Typ,
]

type ExternalKeyer = objekte.ExternalKeyer[
	Akte,
	*Akte,
	kennung.Typ,
	*kennung.Typ,
]

type External = objekte.External[
	Akte,
	*Akte,
	kennung.Typ,
	*kennung.Typ,
]

type CheckedOut = objekte.CheckedOut[
	Akte,
	*Akte,
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
