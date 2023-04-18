package objekte

import (
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type ObjekteLike interface {
	// kennung.Matchable
	metadatei.Getter
}

type ObjektePtrLike interface {
	// kennung.Matchable
	metadatei.Getter
	metadatei.Setter
}
