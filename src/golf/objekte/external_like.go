package objekte

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type ExternalLike interface {
	GetFDs() sku.ExternalFDs
	GetObjekteFD() kennung.FD
	GetAkteFD() kennung.FD
	kennung.Matchable
	metadatei.Getter
}
