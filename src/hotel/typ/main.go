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

func GetFileExtension(t *Transacted) string {
	if t.Akte.FileExtension != "" {
		return t.Akte.FileExtension
	}

	return t.GetKennungLike().String()
}
