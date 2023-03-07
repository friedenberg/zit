package objekte

import "github.com/friedenberg/zit/src/delta/kennung"

type ExternalLike interface {
	GetObjekteFD() kennung.FD
	GetAkteFD() kennung.FD
	kennung.Matchable
}
