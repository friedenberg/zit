package objekte

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type ExternalLike interface {
	GetObjekteFD() kennung.FD
	GetAkteFD() kennung.FD
	kennung.Matchable
	metadatei.Getter
}
